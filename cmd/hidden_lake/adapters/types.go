package adapters

import (
	"context"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IAdaptedConsumer interface {
	Consume(context.Context) (net_message.IMessage, error)
}

type IAdaptedProducer interface {
	Produce(context.Context, net_message.IMessage) error
}
