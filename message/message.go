package message

import (
	"time"

	"github.com/google/uuid"
)

type MessagePrototype struct {
	UUID   *uuid.UUID  `json:"uuid"`
	time   *time.Time  `json:"time"`
	Entity interface{} `json:"entity"`
}

type MessageInterface interface {
	Bytes() []byte
}
