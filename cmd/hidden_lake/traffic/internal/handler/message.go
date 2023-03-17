package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/payload"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func HandleMessageAPI(connKeeper conn_keeper.IConnKeeper, wDB database.IWrapperDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			api.Response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		database := wDB.Get()

		switch r.Method {
		case http.MethodGet:
			query := r.URL.Query()

			msg, err := database.Load(query.Get("hash"))
			if err != nil {
				api.Response(w, pkg_settings.CErrorLoad, "failed: load message")
				return
			}

			api.Response(w, pkg_settings.CErrorNone, encoding.HexEncode(msg.ToBytes()))
			return
		case http.MethodPost:
			var vRequest pkg_settings.SPushRequest

			err := json.NewDecoder(r.Body).Decode(&vRequest)
			if err != nil {
				api.Response(w, pkg_settings.CErrorDecode, "failed: decode request")
				return
			}

			if uint64(len(vRequest.FMessage)/2) > database.Settings().GetMessageSize() {
				api.Response(w, pkg_settings.CErrorPackSize, "failed: incorrect package size")
				return
			}

			msg := message.LoadMessage(
				encoding.HexDecode(vRequest.FMessage),
				message.NewParams(
					database.Settings().GetMessageSize(),
					database.Settings().GetWorkSize(),
				),
			)
			if msg == nil {
				api.Response(w, pkg_settings.CErrorMessage, "failed: decode message")
				return
			}

			if err := database.Push(msg); err != nil {
				api.Response(w, pkg_settings.CErrorPush, "failed: push message")
				return
			}

			connKeeper.GetNetworkNode().BroadcastPayload(
				payload.NewPayload(
					hls_settings.CNetworkMask,
					msg.ToBytes(),
				),
			)

			api.Response(w, pkg_settings.CErrorNone, "success")
			return
		}
	}
}
