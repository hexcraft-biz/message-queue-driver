package message

import (
	"time"

	"github.com/google/uuid"
)

type MessagePrototype struct {
	UUID      *uuid.UUID  `json:"uuid"`
	Timestamp *time.Time  `json:"timestamp"`
	Entity    interface{} `json:"entity"`
}

type MessageInterface interface {
	Bytes() []byte
}
