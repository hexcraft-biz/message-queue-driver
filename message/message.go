package message

import (
	"time"

	"github.com/google/uuid"
)

type MessagePrototype struct {
	UUID   *uuid.UUID  `json:"uuid"`
	Time   *time.Time  `json:"time"`
	Entity interface{} `json:"entity"`
}

func (p *MessagePrototype) Init() {
	uid := uuid.New()
	ctime := time.Now().UTC()

	p.UUID = &uid
	p.Time = &ctime
}

type MessageInterface interface {
	Bytes() []byte
}
