package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/encoding"
)

func HandleLoadAPI(db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vRequest pkg_settings.SLoadRequest

		if r.Method != "POST" {
			response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		err := json.NewDecoder(r.Body).Decode(&vRequest)
		if err != nil {
			response(w, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		msg, err := db.Load(vRequest.FHash)
		if err != nil {
			response(w, pkg_settings.CErrorLoad, "failed: load message")
			return
		}

		response(w, pkg_settings.CErrorNone, encoding.HexEncode(msg.Bytes()))
	}
}
