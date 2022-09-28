package pubsub

import (
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/hexcraft-biz/message-queue-driver/message"
	"golang.org/x/net/context"
)

type PubsubClient struct {
	Entity  *pubsub.Client
	Context context.Context
}

func NewClient(projectID string) (*PubsubClient, error) {
	c := &PubsubClient{}

	c.Context = context.Background()

	client, err := pubsub.NewClient(c.Context, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub: NewClient: %v", err)

	}
	c.Entity = client

	return c, nil
}

func (c *PubsubClient) Close() {
	c.Entity.Close()
}

func (c *PubsubClient) Topic(topicName string) (*PubsubTopic, error) {
	t := &PubsubTopic{
		Entity:  c.Entity.Topic(topicName),
		Context: c.Context,
	}

	if ok, err := t.Entity.Exists(t.Context); err != nil {
		return nil, fmt.Errorf("pubsub: NewClient: %v", err)
	} else if ok == false {
		return nil, fmt.Errorf("pubsub: topic: %s is not exist.", topicName)
	}

	return t, nil
}

type PubsubTopic struct {
	Entity  *pubsub.Topic
	Context context.Context
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
