package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hlm/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hlm/internal/database"
	"github.com/number571/go-peer/cmd/hlm/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	hls_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
)

func HandleIncomigHTTP(db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		msgBytes, err := io.ReadAll(r.Body)
		if err != nil {
			response(w, hls_settings.CErrorResponse, "failed: response message")
			return
		}

		msg := strings.TrimSpace(string(msgBytes))
		if len(msg) == 0 {
			response(w, hls_settings.CErrorResponse, "failed: message is null")
			return
		}

		pubKey := asymmetric.LoadRSAPubKey(r.Header.Get(hls_settings.CHeaderPubKey))
		if pubKey == nil {
			panic("public key is null (receive from hls)!")
		}

		dbMsg := database.NewMessage(true, msg)
		if err := db.Push(pubKey, dbMsg); err != nil {
			response(w, hls_settings.CErrorPubKey, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:   pubKey.Address().String(),
			FMessage:   dbMsg.GetMessage(),
			FTimestamp: dbMsg.GetTimestamp(),
		})
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
