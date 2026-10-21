package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/hexcraft-biz/message-queue-driver/message"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ================================================================
// Client
// ================================================================
func NewClient(projectID string) (*PubsubClient, error) {
	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub: NewClient: %v", err)
	} else if client == nil {
		return nil, fmt.Errorf("pubsub: NewClient: client is null")
	}

	return &PubsubClient{
		Entity: client,
	}, nil
}

type PubsubClient struct {
	Entity *pubsub.Client
}

// Close releases the client connections. Close all topics before closing the client.
func (c *PubsubClient) Close() {
	c.Entity.Close()
}

func (c *PubsubClient) Topic(topicName string) *PubsubTopic {
	return &PubsubTopic{
		Entity: c.Entity.Publisher(topicName),
		client: c.Entity,
	}
}

// ================================================================
// Topic
// ================================================================
// PubsubTopic wraps a v2 publisher. Create topics through PubsubClient.Topic
// so that Exists can use the associated admin client.
type PubsubTopic struct {
	Entity *pubsub.Publisher
	client *pubsub.Client
}

// Close flushes pending messages and stops the publisher without closing the client.
func (t *PubsubTopic) Close() {
	t.Entity.Stop()
}

func (t *PubsubTopic) Exists() (bool, error) {
	if _, err := t.client.TopicAdminClient.GetTopic(context.Background(), &pubsubpb.GetTopicRequest{
		Topic: t.Entity.String(),
	}); status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("pubsub: PubsubTopic.Exists: %w", err)
	}

	return true, nil
}

func (t *PubsubTopic) Publish(msgData message.MessageInterface, attrs map[string]string) *pubsub.PublishResult {
	// Publish message to the topic
	return t.Entity.Publish(context.Background(), &pubsub.Message{
		Data:       msgData.Bytes(),
		Attributes: attrs,
	})
}
