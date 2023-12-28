package queue_pusher

import (
	"sync"

	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/wrapper"
)

var (
	_ IQPWrapper = &sQPWrapper{}
)

type sQPWrapper struct {
	fMutex   sync.Mutex
	fWrapper wrapper.IWrapper
}

func NewQPWrapper() IQPWrapper {
	return &sQPWrapper{fWrapper: wrapper.NewWrapper()}
}

func (p *sQPWrapper) Get() queue_set.IQueuePusher {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	queueSet, ok := p.fWrapper.Get().(queue_set.IQueuePusher)
	if !ok {
		return nil
	}

	return queueSet
}

func (p *sQPWrapper) Set(pQueueSet queue_set.IQueuePusher) IQPWrapper {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fWrapper.Set(pQueueSet)
	return p
}
