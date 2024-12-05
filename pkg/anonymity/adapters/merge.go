package adapters

import (
	"context"
	"errors"
	"sync"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IAdapter = &sMergedAdapter{}
)

type sMergedAdapter struct {
	fMsgChan  chan net_message.IMessage
	fAdapters []*sWaitAdapter
}

type sWaitAdapter struct {
	fMutex   sync.Mutex
	fAdapter IAdapter
}

func NewMergedAdapter(pAdapters ...IAdapter) IAdapter {
	return &sMergedAdapter{
		fMsgChan:  make(chan net_message.IMessage, len(pAdapters)),
		fAdapters: toWaitAdapters(pAdapters),
	}
}

func (p *sMergedAdapter) Produce(pCtx context.Context, pMsg net_message.IMessage) error {
	N := len(p.fAdapters)
	errs := make([]error, N)

	wg := &sync.WaitGroup{}
	wg.Add(N)
	for i, a := range p.fAdapters {
		go func(i int, a *sWaitAdapter) {
			defer wg.Done()
			errs[i] = a.fAdapter.Produce(pCtx, pMsg)
		}(i, a)
	}
	wg.Wait()

	return errors.Join(errs...)
}

func (p *sMergedAdapter) Consume(pCtx context.Context) (net_message.IMessage, error) {
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fMsgChan:
		return msg, nil
	default:
	}

	for _, a := range p.fAdapters {
		go func(a *sWaitAdapter) {
			if ok := a.fMutex.TryLock(); !ok {
				return
			}
			defer a.fMutex.Unlock()
			msg, err := a.fAdapter.Consume(pCtx)
			if err != nil {
				return
			}
			p.fMsgChan <- msg
		}(a)
	}

	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fMsgChan:
		return msg, nil
	}
}

func toWaitAdapters(pAdapters []IAdapter) []*sWaitAdapter {
	result := make([]*sWaitAdapter, 0, len(pAdapters))
	for _, a := range pAdapters {
		result = append(result, &sWaitAdapter{fAdapter: a})
	}
	return result
}
