# message-queue-driver

## Quick start

### Assume the following below codes

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hexcraft-biz/message-queue-driver/pubsub"
)

type Msg struct {
	ID string `json:"id"`
}

func (m *Msg) Bytes() []byte {
	msgData, _ := json.Marshal(m)

	return msgData
}

type Text struct {
	Content string
}

func (t *Text) Bytes() []byte {
	return []byte(t.Content)
}

func main() {
	// export GOOGLE_APPLICATION_CREDENTIALS=/yourpath/credentials.jso
	projectID := "projectID"
	topicName := "topicName"

	c, err := pubsub.NewPubsubClient(projectID)
	fmt.Println(c, err)

	t, err := c.Topic(topicName)
	fmt.Println(t, err)

	msgId, err := t.Publish(&Msg{
		ID: "test-123123",
	})
	fmt.Println(msgId, err)

	textId, err := t.Publish(&Text{
		Content: "Hello World",
	})
	fmt.Println(textId, err)

	c.Close()
}
```
