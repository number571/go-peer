package adapters

import (
	"context"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IAdapter = &sAdapter{}
)

type (
	iProducerF func(context.Context, net_message.IMessage) error
	iConsumerF func(context.Context) (net_message.IMessage, error)
)

type sAdapter struct {
	fProduce iProducerF
	fConsume iConsumerF
}

func NewAdapterByFuncs(pProduce iProducerF, pConsume iConsumerF) IAdapter {
	return &sAdapter{
		fProduce: pProduce,
		fConsume: pConsume,
	}
}

func (p *sAdapter) Produce(pCtx context.Context, pMsg net_message.IMessage) error {
	return p.fProduce(pCtx, pMsg)
}

func (p *sAdapter) Consume(pCtx context.Context) (net_message.IMessage, error) {
	return p.fConsume(pCtx)
}
