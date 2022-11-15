package handler

import (
	"sync"

	"golang.org/x/net/websocket"
)

type sChatWS struct {
	FAddress string `json:"address"`
	FMessage string `json:"message"`
}

var (
	gPrevConnWS *websocket.Conn
	gMutexWS    sync.Mutex
)

var (
	gSignalWS = make(chan struct{})
	gChatWS   = make(chan *sChatWS)
)

func FriendsChatWS(ws *websocket.Conn) {
	defer ws.Close()

	getIncoming := new(sChatWS)
	if err := websocket.JSON.Receive(ws, getIncoming); err != nil {
		return
	}

	gMutexWS.Lock()
	if gPrevConnWS != nil {
		gSignalWS <- struct{}{}
	}
	gPrevConnWS = ws
	gMutexWS.Unlock()

	for {
		select {
		case incoming := <-gChatWS:
			if getIncoming.FAddress != incoming.FAddress {
				continue
			}
			if err := websocket.JSON.Send(ws, incoming); err != nil {
				gPrevConnWS = nil
				return
			}
		case <-gSignalWS:
			return
		}
	}
}
