package handler

import (
	"encoding/json"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkPushAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vPush pkg_settings.SPush

		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			api.Response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vPush); err != nil {
			api.Response(w, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		pubKey := asymmetric.LoadRSAPubKey(vPush.FReceiver)
		if pubKey == nil {
			api.Response(w, pkg_settings.CErrorPubKey, "failed: load public key")
			return
		}

		data := encoding.HexDecode(vPush.FHexData)
		if data == nil {
			api.Response(w, pkg_settings.CErrorPubKey, "failed: decode hex format data")
			return
		}

		switch r.Method {
		case http.MethodPut:
			err := node.Broadcast(
				pubKey,
				anonymity.NewPayload(pkg_settings.CHeaderHLS, data),
			)
			if err != nil {
				api.Response(w, pkg_settings.CErrorMessage, "failed: broadcast message")
				return
			}
			api.Response(w, pkg_settings.CErrorNone, "success: broadcast")
			return
		case http.MethodPost:
			resp, err := node.Request(
				pubKey,
				anonymity.NewPayload(pkg_settings.CHeaderHLS, data),
			)
			if err != nil {
				api.Response(w, pkg_settings.CErrorResponse, "failed: response message")
				return
			}
			api.Response(w, pkg_settings.CErrorNone, encoding.HexEncode(resp))
			return
		}
	}
}
