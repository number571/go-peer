package handler

import (
	"encoding/json"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkMessageAPI(pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		var vMessage pkg_settings.SMessage

		if pR.Method != http.MethodPost {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vMessage); err != nil {
			api.Response(pW, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		msg := message.LoadMessage(
			pNode.GetMessageQueue().GetClient().GetSettings(),
			encoding.HexDecode(vMessage.FHexMessage),
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
		// in future: may be decrypt functions
	}
}
