package msgbroker

import (
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
)

var (
	_ IMessageBroker = &sMessageBroker{}
)

type sMessageBroker struct {
	fQueue   chan sSubscribeMessage
	fCancel  chan struct{}
	fMutex   sync.Mutex
	fConsume bool
}

type sSubscribeMessage struct {
	utils.SSubscribe
	utils.SMessage
}

func NewMessageBroker() IMessageBroker {
	return &sMessageBroker{
		fQueue:  make(chan sSubscribeMessage, 1),
		fCancel: make(chan struct{}),
	}
}

func (p *sMessageBroker) Consume(pAddress string) (utils.SMessage, bool) {
	p.fMutex.Lock()
	if p.fConsume {
		p.fCancel <- struct{}{}
	}
	p.fConsume = true
	p.fMutex.Unlock()

	select {
	case msg, ok := <-p.fQueue:
		p.fMutex.Lock()
		p.fConsume = false
		p.fMutex.Unlock()
		return msg.SMessage, ok && msg.FAddress == pAddress
	case <-p.fCancel:
		// no need set consume = false,
		// one consumer close another consumer
		return utils.SMessage{}, false
	}
}

func (p *sMessageBroker) Produce(pAddress string, pMsg utils.SMessage) {
	// clear the queue if there are no consumers
	for len(p.fQueue) > 0 {
		<-p.fQueue
	}
	p.fQueue <- sSubscribeMessage{
		SSubscribe: utils.SSubscribe{FAddress: pAddress},
		SMessage:   pMsg,
	}
}
