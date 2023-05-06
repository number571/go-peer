package handler

import (
	"io"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkMessageAPI(pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodPost {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		// get message from HLT
		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			api.Response(pW, pkg_settings.CErrorRead, "failed: read message bytes")
			return
		}

		msg := message.LoadMessage(
			pNode.GetMessageQueue().GetClient().GetSettings(),
			msgBytes,
		)
		if msg == nil {
			api.Response(pW, pkg_settings.CErrorMessage, "failed: decode hex format message")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			pNode.HandleMessage(msg)
			api.Response(pW, pkg_settings.CErrorNone, "success: handle message")
			return
		}
		// may be decrypt functions
	}
}
