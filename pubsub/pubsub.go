package pubsub

import (
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/hexcraft-biz/message-queue-driver/message"
	"golang.org/x/net/context"
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

func (c *PubsubClient) Close() {
	c.Entity.Close()
}

func (c *PubsubClient) Topic(topicName string) *PubsubTopic {
	return &PubsubTopic{
		Entity: c.Entity.Topic(topicName),
	}
}

// ================================================================
// Topic
// ================================================================
type PubsubTopic struct {
	Entity *pubsub.Topic
}

func (t *PubsubTopic) Exists() (bool, error) {
	if ok, err := t.Entity.Exists(context.Background()); err != nil {
		return false, fmt.Errorf("pubsub: PubsubTopic.Exists: %v", err)
	} else {
		return ok, nil
	}
}

func (t *PubsubTopic) Publish(msgData message.MessageInterface, attrs map[string]string) *pubsub.PublishResult {
	defer t.Entity.Stop()

	// Publish message to the topic
	res := t.Entity.Publish(context.Background(), &pubsub.Message{
		Data:       msgData.Bytes(),
		Attributes: attrs,
	})

	return res
}
