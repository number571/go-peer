package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
)

func HandlePushAPI(db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vRequest pkg_settings.SPushRequest

		if r.Method != "POST" {
			response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		err := json.NewDecoder(r.Body).Decode(&vRequest)
		if err != nil {
			response(w, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		if uint64(len(vRequest.FMessage)) > hlt_settings.CMessageSize {
			response(w, pkg_settings.CErrorPackSize, "failed: incorrect package size")
			return
		}

		msg := message.LoadMessage(encoding.HexDecode(vRequest.FMessage), hlt_settings.CWorkSize)
		if msg == nil {
			response(w, pkg_settings.CErrorMessage, "failed: decode message")
			return
		}

		err = db.Push(msg)
		if err != nil {
			response(w, pkg_settings.CErrorPush, "failed: push message")
			return
		}

		response(w, pkg_settings.CErrorNone, "success")
	}
}
