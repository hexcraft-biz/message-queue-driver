package pubsub

import (
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/hexcraft-biz/message-queue-driver/message"
	"golang.org/x/net/context"
)

func NewClient(projectID string) (*PubsubClient, error) {
	context := context.Background()

	if client, err := pubsub.NewClient(context, projectID); err != nil {
		return nil, fmt.Errorf("pubsub: NewClient: %v", err)
	} else {
		return &PubsubClient{
			Entity:  client,
			Context: context,
		}, nil
	}
}

type PubsubClient struct {
	Entity  *pubsub.Client
	Context context.Context
}

func (c *PubsubClient) Close() {
	c.Entity.Close()
}

func (c *PubsubClient) Topic(topicName string) (*PubsubTopic, error) {
	return &PubsubTopic{
		Entity:  c.Entity.Topic(topicName),
		Context: c.Context,
	}, nil
}

type PubsubTopic struct {
	Entity  *pubsub.Topic
	Context context.Context
}

func (t *PubsubTopic) Exists() (bool, error) {
	if ok, err := t.Entity.Exists(t.Context); err != nil {
		return false, fmt.Errorf("pubsub: NewClient: %v", err)
	} else {
		return ok, nil
	}
}

func (t *PubsubTopic) Publish(m message.MessageInterface) (string, error) {
	res := t.Entity.Publish(t.Context, &pubsub.Message{
		Data: m.Bytes(),
	})

	msgId, err := res.Get(t.Context)
	if err != nil {
		return "", fmt.Errorf("pubsub: result.Get: %v", err)
	}

	return msgId, nil
}
