package queue

import (
	"time"

	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/client/message"
)

type IQueue interface {
	Settings() ISettings
	Client() client.IClient

	Enqueue(message.IMessage) error
	Dequeue() <-chan message.IMessage

	modules.IApp
}

type ISettings interface {
	GetCapacity() uint64
	GetPullCapacity() uint64
	GetDuration() time.Duration
}
