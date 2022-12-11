package handler

import (
	"github.com/number571/go-peer/cmd/hlm/chat_queue"
	"golang.org/x/net/websocket"
)

var (
	gChatQueue = chat_queue.NewChatQueue(1)
)

func FriendsChatWS(ws *websocket.Conn) {
	defer ws.Close()

	subscribe := new(chat_queue.SMessage)
	if err := websocket.JSON.Receive(ws, subscribe); err != nil {
		return
	}

	gChatQueue.Init()
	for {
		msg, ok := gChatQueue.Load(subscribe.FAddress)
		if !ok {
			return
		}
		if err := websocket.JSON.Send(ws, msg); err != nil {
			return
		}
	}
}
