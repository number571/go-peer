package adapters

import (
	"context"

	"github.com/number571/go-peer/pkg/message/layer1"
)

var (
	_ IAdapter = &sAdapter{}
)

type (
	iProducerF func(context.Context, layer1.IMessage) error
	iConsumerF func(context.Context) (layer1.IMessage, error)
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

func (p *sAdapter) Produce(pCtx context.Context, pMsg layer1.IMessage) error {
	return p.fProduce(pCtx, pMsg)
}

func (p *sAdapter) Consume(pCtx context.Context) (layer1.IMessage, error) {
	return p.fConsume(pCtx)
}
