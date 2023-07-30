package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkMessageAPI(pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodPost {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			api.Response(pW, http.StatusConflict, "failed: read message bytes")
			return
		}

		msg := message.LoadMessage(
			pNode.GetMessageQueue().GetClient().GetSettings(),
			string(msgBytes),
		)
		if msg == nil {
			api.Response(pW, http.StatusBadRequest, "failed: decode message")
			return
		}

		pNode.HandleMessage(msg)
		api.Response(pW, http.StatusOK, "success: handle message")
	}
}
