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

func NewChatQueue(size int) IChatQueue {
	return &sChatQueue{
		fQueue: make(chan *SMessage, size),
		fSize:  size,
	}
}

func (cq *sChatQueue) Init() {
	defer time.Sleep(200 * time.Millisecond)

	cq.fMutex.Lock()
	defer cq.fMutex.Unlock()

	close(cq.fQueue)
	cq.fQueue = make(chan *SMessage, cq.fSize)
}

func (cq *sChatQueue) Load(addr string) (*SMessage, bool) {
	defer func() {
		cq.fMutex.Lock()
		cq.fListen = ""
		cq.fMutex.Unlock()
	}()

	cq.fMutex.Lock()
	cq.fListen = addr
	cq.fMutex.Unlock()

	msg, ok := <-cq.fQueue
	return msg, ok
}

func (cq *sChatQueue) Push(msg *SMessage) {
	cq.fMutex.Lock()
	defer cq.fMutex.Unlock()

	if cq.fListen != msg.FAddress {
		return
	}

	// to prevent blocking response
	if len(cq.fQueue) == cq.fSize {
		return
	}

	cq.fQueue <- msg
}
