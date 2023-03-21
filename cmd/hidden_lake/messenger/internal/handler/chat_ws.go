package handler

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/chat_queue"
	"golang.org/x/net/websocket"
)

var (
	gChatQueue = chat_queue.NewChatQueue(1)
)

func FriendsChatWS(pWS *websocket.Conn) {
	defer pWS.Close()

	subscribe := new(chat_queue.SMessage)
	if err := websocket.JSON.Receive(pWS, subscribe); err != nil {
		return
	}

	gChatQueue.Init()
	for {
		msg, ok := gChatQueue.Load(subscribe.FAddress)
		if !ok {
			return
		}
		if err := websocket.JSON.Send(pWS, msg); err != nil {
			return
		}
	}
}
