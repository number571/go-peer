package adapters

import (
	"context"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IAdapter interface {
	IProducer
	IConsumer
}

type IProducer interface {
	Produce(context.Context, net_message.IMessage) error
}

type IConsumer interface {
	Consume(context.Context) (net_message.IMessage, error)
}
