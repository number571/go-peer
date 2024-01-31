package receiver

import (
	"sync"
)

var (
	_ IMessageReceiver = &sMessageReceiver{}
)

type sMessageReceiver struct {
	fMutex sync.Mutex
	fQueue chan *SMessage
}

func NewMessageReceiver() IMessageReceiver {
	return &sMessageReceiver{
		fQueue: make(chan *SMessage),
	}
}

func (p *sMessageReceiver) Init() IMessageReceiver {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	close(p.fQueue)
	p.fQueue = make(chan *SMessage)
	return p
}

func (p *sMessageReceiver) Recv(pAddress string) (*SMessage, bool) {
	msg, ok := <-p.fQueue

	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return msg, ok && msg.FAddress == pAddress
}

func (p *sMessageReceiver) Send(msg *SMessage) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fQueue <- msg
}
