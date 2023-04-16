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

func HandleMessageAPI(pConnKeeper conn_keeper.IConnKeeper, pWrapperDB database.IWrapperDB) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()

		switch pR.Method {
		case http.MethodGet:
			query := pR.URL.Query()

			msg, err := database.Load(query.Get("hash"))
			if err != nil {
				api.Response(pW, pkg_settings.CErrorLoad, "failed: load message")
				return
			}

			api.Response(pW, pkg_settings.CErrorNone, encoding.HexEncode(msg.ToBytes()))
			return
		case http.MethodPost:
			var vRequest pkg_settings.SPushRequest

			err := json.NewDecoder(pR.Body).Decode(&vRequest)
			if err != nil {
				api.Response(pW, pkg_settings.CErrorDecode, "failed: decode request")
				return
			}

			if uint64(len(vRequest.FMessage)/2) > database.Settings().GetMessageSize() {
				api.Response(pW, pkg_settings.CErrorPackSize, "failed: incorrect package size")
				return
			}

			msg := message.LoadMessage(
				message.NewSettings(&message.SSettings{
					FMessageSize: database.Settings().GetMessageSize(),
					FWorkSize:    database.Settings().GetWorkSize(),
				}),
				encoding.HexDecode(vRequest.FMessage),
			)
			if msg == nil {
				api.Response(pW, pkg_settings.CErrorMessage, "failed: decode message")
				return
			}

			if err := database.Push(msg); err != nil {
				api.Response(pW, pkg_settings.CErrorPush, "failed: push message")
				return
			}

			pConnKeeper.GetNetworkNode().BroadcastPayload(
				payload.NewPayload(
					hls_settings.CNetworkMask,
					msg.ToBytes(),
				),
			)

			api.Response(pW, pkg_settings.CErrorNone, "success")
			return
		}
	}
}
