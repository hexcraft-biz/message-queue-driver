package pubsub

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type testPublisherServer struct {
	pubsubpb.UnimplementedPublisherServer
	published chan *pubsubpb.PublishRequest
}

func (s *testPublisherServer) GetTopic(_ context.Context, req *pubsubpb.GetTopicRequest) (*pubsubpb.Topic, error) {
	switch req.Topic {
	case "projects/test-project/topics/existing":
		return &pubsubpb.Topic{Name: req.Topic}, nil
	case "projects/test-project/topics/missing":
		return nil, status.Error(codes.NotFound, "missing topic")
	default:
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}
}

func (s *testPublisherServer) Publish(_ context.Context, req *pubsubpb.PublishRequest) (*pubsubpb.PublishResponse, error) {
	s.published <- req
	return &pubsubpb.PublishResponse{MessageIds: []string{"message-id"}}, nil
}

func newTestClient(t *testing.T) (*PubsubClient, *testPublisherServer) {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	backend := &testPublisherServer{published: make(chan *pubsubpb.PublishRequest, 4)}
	pubsubpb.RegisterPublisherServer(server, backend)
	go server.Serve(listener)
	t.Cleanup(func() {
		server.Stop()
		listener.Close()
	})

	conn, err := grpc.NewClient("passthrough:///pubsub-test",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	entity, err := pubsub.NewClient(context.Background(), "test-project", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatal(err)
	}
	client := &PubsubClient{Entity: entity}
	t.Cleanup(client.Close)
	return client, backend
}

func TestTopicExists(t *testing.T) {
	client, _ := newTestClient(t)
	for _, tc := range []struct {
		name string
		want bool
		code codes.Code
	}{
		{"existing", true, codes.OK},
		{"projects/test-project/topics/existing", true, codes.OK},
		{"missing", false, codes.OK},
		{"forbidden", false, codes.PermissionDenied},
	} {
		t.Run(tc.name, func(t *testing.T) {
			topic := client.Topic(tc.name)
			defer topic.Close()
			got, err := topic.Exists()
			if got != tc.want || status.Code(err) != tc.code {
				t.Fatalf("Exists() = (%v, %v), want (%v, %v)", got, err, tc.want, tc.code)
			}
		})
	}
}

type testMessage []byte

func (m testMessage) Bytes() []byte { return m }

func TestPublishAndTopicClose(t *testing.T) {
	client, backend := newTestClient(t)
	topic := client.Topic("existing")
	defer topic.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := testMessage(`{"uuid":"id","time":"2026-10-21T00:00:00Z","entity":{"text":"hello"}}`)
	attrs := map[string]string{"tracking_id": "tracking-id"}
	result := topic.Publish(data, attrs)
	// Closing must flush queued messages before the shared client is closed.
	topic.Close()
	id, err := result.Get(ctx)
	if err != nil || id != "message-id" {
		t.Fatalf("Get() = (%q, %v)", id, err)
	}
	select {
	case req := <-backend.published:
		if req.Topic != "projects/test-project/topics/existing" || len(req.Messages) != 1 {
			t.Fatalf("unexpected publish request: %v", req)
		}
		msg := req.Messages[0]
		if string(msg.Data) != string(data) || !reflect.DeepEqual(msg.Attributes, attrs) {
			t.Fatalf("payload or attributes changed: %v", msg)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if _, err := topic.Publish(data, nil).Get(ctx); !errors.Is(err, pubsub.ErrPublisherStopped) {
		t.Fatalf("publish after Close() error = %v", err)
	}
	// A closed publisher must not close the shared client or another publisher.
	other := client.Topic("other")
	defer other.Close()
	if _, err := other.Publish(data, nil).Get(ctx); err != nil {
		t.Fatalf("publish through shared client: %v", err)
	}
}
