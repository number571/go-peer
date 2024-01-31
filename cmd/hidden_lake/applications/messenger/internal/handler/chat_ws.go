package handler

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/receiver"
	"golang.org/x/net/websocket"
)

func FriendsChatWS(pReceiver receiver.IMessageReceiver) func(pWS *websocket.Conn) {
	return func(pWS *websocket.Conn) {
		defer pWS.Close()

		subscribe := new(receiver.SMessage)
		if err := websocket.JSON.Receive(pWS, subscribe); err != nil {
			return
		}

		pReceiver.Init()
		for {
			msg, ok := pReceiver.Recv(subscribe.FAddress)
			if !ok {
				return
			}
			if err := websocket.JSON.Send(pWS, msg); err != nil {
				return
			}
		}
	}
}
