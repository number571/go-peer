package queue

import (
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/types"
)

type IMessageQueue interface {
	types.ICommand

	GetSettings() ISettings
	GetClient() client.IClient

	UpdateClient(client.IClient)

	EnqueueMessage(message.IMessage) error
	DequeueMessage() <-chan message.IMessage
}

type ISettings interface {
	GetCapacity() uint64
	GetPullCapacity() uint64
	GetDuration() time.Duration
}
