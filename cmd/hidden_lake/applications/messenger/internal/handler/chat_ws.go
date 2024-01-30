package handler

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/receiver"
	"golang.org/x/net/websocket"
)

var (
	gReceiver = receiver.NewMessageReceiver()
)

func FriendsChatWS(pWS *websocket.Conn) {
	defer pWS.Close()

	subscribe := new(receiver.SMessage)
	if err := websocket.JSON.Receive(pWS, subscribe); err != nil {
		return
	}

	gReceiver.Init(subscribe.FAddress)
	for {
		msg, ok := gReceiver.Recv()
		if !ok {
			return
		}
		if err := websocket.JSON.Send(pWS, msg); err != nil {
			return
		}
	}
}
