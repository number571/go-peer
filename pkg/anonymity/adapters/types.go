package adapters

import (
	"context"

	"github.com/number571/go-peer/pkg/message/layer1"
)

type IAdapter interface {
	IProducer
	IConsumer
}

type IProducer interface {
	Produce(context.Context, layer1.IMessage) error
}

type IConsumer interface {
	Consume(context.Context) (layer1.IMessage, error)
}
