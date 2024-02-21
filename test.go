package main

import (
	"encoding/json"

	"github.com/hexcraft-biz/message-queue-driver/message"
	// "github.com/mitchellh/mapstructure"
)

type Msg struct {
	message.MessagePrototype
}

type mdata struct {
	Text string `json:"text"`
}

func (m *Msg) Bytes() []byte {
	msgData, _ := json.Marshal(m)

	return msgData
}

func main() {

	// jj := `{"uuid":"353a1e3f-45bc-4c25-bc5e-5f3867297852","time":"2022-09-25T07:57:30Z","entity":{"aaa":"AAA","bbb":"BBB","text":"test-123123"}}`

	/*
		entity, mm := mdata{}, Msg{}
		mm.Entity = &entity
		json.Unmarshal([]byte(jj), &mm)
		fmt.Println(mm)
		fmt.Println(mm.Entity.(*mdata).Text)
	*/

	/*
		entity, mm := mdata{}, Msg{}
		json.Unmarshal([]byte(jj), &mm)
		fmt.Println(mm)
		mapstructure.Decode(mm.Entity, &entity)
		fmt.Println(entity, entity.Text)
	*/

	// export GOOGLE_APPLICATION_CREDENTIALS=/yourpath/credentials.json
	/*
		projectID := "dev-20221001-chatisfy"
		topicName := "john-test-123"
		c, err := pubsub.NewClient(projectID)
		fmt.Println(c, err)

		t, err := c.Topic(topicName)
		fmt.Println(t, err)

		uid := uuid.New()
		ctime := time.Now().UTC().Truncate(time.Second)

		for i := 0; i < 100; i++ {
			msgId, err := t.Publish(&Msg{
				MessagePrototype: message.MessagePrototype{
					UUID: &uid,
					Time: &ctime,
					Entity: mdata{
						Text: "test-" + strconv.Itoa(i),
					},
				},
			})
			fmt.Println(msgId, err)
		}
	*/

	/*
		msgId1, err := t.Publish(&Msg{
			MessagePrototype: message.MessagePrototype{
				UUID:   &uid,
				Time:   &ctime,
				Entity: "ok",
			},
		})
		fmt.Println(msgId1, err)
	*/

	// c.Close()
}
