package queue

import (
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/message"
)

type IQueue interface {
	Settings() ISettings
	Client() client.IClient

	Run() error
	Close() error

	Enqueue(message.IMessage) error
	Dequeue() <-chan message.IMessage
}

type ISettings interface {
	GetCapacity() uint64
	GetPullCapacity() uint64
	GetDuration() time.Duration
}
