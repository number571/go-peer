package handler

import (
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hlm/database"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"golang.org/x/net/websocket"
)

type sMessage struct {
	FAddress   string `json:"address"`
	FMessage   string `json:"message"`
	FTimestamp string `json:"timestamp"`
}

type sChatQueue struct {
	fQueue  chan *sMessage
	fMutex  sync.Mutex
	fListen string
}

const (
	queueSize = 1
)

var (
	gChatQueue = newChatQueue()
)

func FriendsChatWS(ws *websocket.Conn) {
	defer ws.Close()

	subscribe := new(sMessage)
	if err := websocket.JSON.Receive(ws, subscribe); err != nil {
		return
	}

	gChatQueue.initQueue()
	time.Sleep(200 * time.Millisecond)

	for {
		msg, ok := gChatQueue.loadMessage(subscribe.FAddress)
		if !ok {
			return
		}
		if err := websocket.JSON.Send(ws, msg); err != nil {
			return
		}
	}
}

func newChatQueue() *sChatQueue {
	return &sChatQueue{
		fQueue: make(chan *sMessage, queueSize),
	}
}

func (cq *sChatQueue) initQueue() {
	cq.fMutex.Lock()
	defer cq.fMutex.Unlock()

	close(cq.fQueue)
	cq.fQueue = make(chan *sMessage, queueSize)
}

func (cq *sChatQueue) loadMessage(addr string) (*sMessage, bool) {
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

func (cq *sChatQueue) pushMessage(pubKey asymmetric.IPubKey, msg database.IMessage) {
	cq.fMutex.Lock()
	defer cq.fMutex.Unlock()

	if cq.fListen != pubKey.Address().String() {
		return
	}

	// to prevent blocking response
	if len(cq.fQueue) == queueSize {
		return
	}

	cq.fQueue <- &sMessage{
		FAddress:   pubKey.Address().String(),
		FMessage:   msg.GetMessage(),
		FTimestamp: msg.GetTimestamp(),
	}
}
