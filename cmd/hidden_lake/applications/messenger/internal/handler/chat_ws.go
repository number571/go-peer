package handler

import (
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/msgbroker"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"golang.org/x/net/websocket"
)

func FriendsChatWS(pBroker msgbroker.IMessageBroker) func(pWS *websocket.Conn) {
	return func(pWS *websocket.Conn) {
		defer pWS.Close()

		subscribe := new(utils.SSubscribe)
		if err := websocket.JSON.Receive(pWS, subscribe); err != nil {
			return
		}

		for {
			msg, ok := pBroker.Consume(subscribe.FAddress)
			if !ok {
				return
			}
			if err := websocket.JSON.Send(pWS, msg); err != nil {
				return
			}
		}
	}
}
