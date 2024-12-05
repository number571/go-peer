package anonymity

import (
	"context"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type sAdapter struct {
	fProduce IProducerF
	fConsume IConsumerF
}

func NewAdapterByFuncs(pProduce IProducerF, pConsume IConsumerF) IAdapter {
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
