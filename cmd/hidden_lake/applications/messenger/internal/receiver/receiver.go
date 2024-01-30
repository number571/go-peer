package receiver

import "sync"

var (
	_ IMessageReceiver = &sMessageReceiver{}
)

type sMessageReceiver struct {
	fMutex   sync.Mutex
	fQueue   chan *SMessage
	fAddress string
}

func NewMessageReceiver() IMessageReceiver {
	return &sMessageReceiver{
		fQueue: make(chan *SMessage),
	}
}

func (p *sMessageReceiver) Init(pAddress string) IMessageReceiver {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fAddress = pAddress
	return p
}

func (p *sMessageReceiver) Recv() (*SMessage, bool) {
	msg, ok := <-p.fQueue
	if !ok {
		return nil, false
	}

	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if msg.FAddress != p.fAddress {
		return nil, false
	}

	return msg, true
}

func (p *sMessageReceiver) Send(msg *SMessage) {
	p.fQueue <- msg
}
