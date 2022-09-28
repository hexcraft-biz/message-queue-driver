# message-queue-driver

## Quick start

### Assume the following below codes

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hexcraft-biz/message-queue-driver/message"
	"github.com/hexcraft-biz/message-queue-driver/pubsub"
)

type Msg struct {
	message.MessagePrototype
}

type mdata struct {
	A    string `json:"aaa"`
	B    string `json:"bbb"`
	Text string `json:"text"`
}

func (m *Msg) Bytes() []byte {
	msgData, _ := json.Marshal(m)

	return msgData
}

func main() {

	/*
		jj := `{"uuid":"353a1e3f-45bc-4c25-bc5e-5f3867297852","timestamp":"2022-09-25T07:57:30Z","entity":{"aaa":"AAA","bbb":"BBB","text":"test-123123"}}`
		mm := msg.MessagePrototype{}
		json.Unmarshal([]byte(jj), &mm)
		fmt.Println(mm)
	*/

	// export GOOGLE_APPLICATION_CREDENTIALS=/yourpath/credentials.json
	projectID := "dev-20220901-chatisfy"
	topicName := "test-test-123"
	c, err := pubsub.NewClient(projectID)
	fmt.Println(c, err)

	t, err := c.Topic(topicName)
	fmt.Println(t, err)

	uid := uuid.New()
	ctime := time.Now().UTC().Truncate(time.Second)

	msgId, err := t.Publish(&Msg{
		MessagePrototype: message.MessagePrototype{
			UUID:      &uid,
			Timestamp: &ctime,
			Entity: mdata{
				A:    "AAA",
				B:    "BBB",
				Text: "test-123123",
			},
		},
	})
	fmt.Println(msgId, err)

	msgId1, err := t.Publish(&Msg{
		MessagePrototype: message.MessagePrototype{
			UUID:      &uid,
			Timestamp: &ctime,
			Entity:    "ok",
		},
	})
	fmt.Println(msgId1, err)

	c.Close()
}
```
