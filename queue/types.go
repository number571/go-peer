package queue

import (
	"time"

	"github.com/number571/go-peer/message"
)

type IQueue interface {
	Settings() ISettings

	Start() error
	Close() error

	Enqueue(message.IMessage) error
	Dequeue() <-chan message.IMessage
}

type ISettings interface {
	GetMainCapacity() uint64
	GetPullCapacity() uint64
	GetDuration() time.Duration
}
