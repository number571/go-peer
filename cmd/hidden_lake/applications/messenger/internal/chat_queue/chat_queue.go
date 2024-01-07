package chat_queue

import (
	"sync"
	"time"
)

var (
	_ IChatQueue = &sChatQueue{}
)

type sChatQueue struct {
	fQueue  chan *SMessage
	fMutex  sync.Mutex
	fListen string
	fSize   int
}

func NewChatQueue(pSize int) IChatQueue {
	return &sChatQueue{
		fQueue: make(chan *SMessage, pSize),
		fSize:  pSize,
	}
}

func (p *sChatQueue) Init() {
	defer time.Sleep(200 * time.Millisecond)

	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	close(p.fQueue)
	p.fQueue = make(chan *SMessage, p.fSize)
}

func (p *sChatQueue) Load(addr string) (*SMessage, bool) {
	defer func() {
		p.fMutex.Lock()
		p.fListen = ""
		p.fMutex.Unlock()
	}()

	p.fMutex.Lock()
	p.fListen = addr
	p.fMutex.Unlock()

	msg, ok := <-p.fQueue
	return msg, ok
}

func (p *sChatQueue) Push(msg *SMessage) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fListen != msg.FAddress {
		return
	}

	// to prevent blocking response
	if len(p.fQueue) == p.fSize {
		return
	}

	p.fQueue <- msg
}
