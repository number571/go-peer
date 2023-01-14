package queue

import (
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"
)

type IQueue interface {
	types.IApp

	Settings() ISettings

	UpdateClient(client.IClient)
	Client() client.IClient

	Enqueue(message.IMessage) error
	Dequeue() <-chan message.IMessage
}

type ISettings interface {
	GetCapacity() uint64
	GetPullCapacity() uint64
	GetDuration() time.Duration
}
