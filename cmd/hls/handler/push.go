package handler

import (
	"encoding/json"
	"net/http"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/network/anonymity"
	payload_adapter "github.com/number571/go-peer/modules/network/anonymity/adapters/payload"
)

func HandlePushAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vPush hls_settings.SPush

		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vPush); err != nil {
			response(w, hls_settings.CErrorDecode, "failed: decode request")
			return
		}

		pubKey := asymmetric.LoadRSAPubKey(vPush.FReceiver)
		if pubKey == nil {
			response(w, hls_settings.CErrorPubKey, "failed: load public key")
			return
		}

		data := encoding.HexDecode(vPush.FHexData)
		if data == nil {
			response(w, hls_settings.CErrorPubKey, "failed: decode hex format data")
			return
		}

		switch r.Method {
		case http.MethodPut:
			msg, err := node.Queue().Client().Encrypt(
				pubKey,
				payload_adapter.NewPayload(hls_settings.CHeaderHLS, data),
			)
			if err != nil {
				response(w, hls_settings.CErrorMessage, "failed: encrypt message with data")
				return
			}
			if err := node.Queue().Enqueue(msg); err != nil {
				response(w, hls_settings.CErrorBroadcast, "failed: broadcast message")
				return
			}
			response(w, hls_settings.CErrorNone, "success: broadcast")
			return
		case http.MethodPost:
			resp, err := node.Request(
				pubKey,
				payload_adapter.NewPayload(hls_settings.CHeaderHLS, data),
			)
			if err != nil {
				response(w, hls_settings.CErrorResponse, "failed: response message")
				return
			}
			response(w, hls_settings.CErrorNone, encoding.HexEncode(resp))
			return
		}
	}
}
