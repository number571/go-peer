package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/database"
	"github.com/number571/go-peer/cmd/hlm/settings"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules/crypto/asymmetric"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func HandleIncomigHTTP(client hlc.IClient, db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		msg, err := io.ReadAll(r.Body)
		if err != nil {
			response(w, hls_settings.CErrorResponse, "failed: response message")
			return
		}

		pubKeyStr := r.Header.Get(hls_settings.CHeaderPubKey)
		err = db.Push(
			asymmetric.LoadRSAPubKey(pubKeyStr),
			database.NewMessage(true, string(msg)),
		)
		if err != nil {
			response(w, hls_settings.CErrorPubKey, "failed: push message to database")
			return
		}

		response(w, hls_settings.CErrorNone, settings.CTitlePattern)
	}
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		FResult: res,
		FReturn: ret,
	})
}
