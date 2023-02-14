package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
)

func HandleMessageAPI(db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			api.Response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		switch r.Method {
		case http.MethodGet:
			query := r.URL.Query()
			msg, err := db.Load(query.Get("hash"))
			if err != nil {
				api.Response(w, pkg_settings.CErrorLoad, "failed: load message")
				return
			}

			api.Response(w, pkg_settings.CErrorNone, encoding.HexEncode(msg.Bytes()))
			return
		case http.MethodPost:
			var vRequest pkg_settings.SPushRequest

			err := json.NewDecoder(r.Body).Decode(&vRequest)
			if err != nil {
				api.Response(w, pkg_settings.CErrorDecode, "failed: decode request")
				return
			}

			if uint64(len(vRequest.FMessage)/2) > db.Settings().GetMessageSize() {
				api.Response(w, pkg_settings.CErrorPackSize, "failed: incorrect package size")
				return
			}

			msg := message.LoadMessage(
				encoding.HexDecode(vRequest.FMessage),
				message.NewParams(
					db.Settings().GetMessageSize(),
					db.Settings().GetWorkSize(),
				),
			)
			if msg == nil {
				api.Response(w, pkg_settings.CErrorMessage, "failed: decode message")
				return
			}

			err = db.Push(msg)
			if err != nil {
				api.Response(w, pkg_settings.CErrorPush, "failed: push message")
				return
			}

			api.Response(w, pkg_settings.CErrorNone, "success")
			return
		}
	}
}
