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

	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
)

func HandleIncomigHTTP(wDB database.IWrapperDB, client hls_client.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		db := wDB.Get()
		if db == nil {
			response(w, hls_settings.CErrorUnauth, "failed: client unauthorized")
			return
		}

		msgBytes, err := io.ReadAll(r.Body)
		if err != nil {
			response(w, hls_settings.CErrorResponse, "failed: response message")
			return
		}

		msg := strings.TrimSpace(string(msgBytes))
		if len(msg) == 0 {
			response(w, hls_settings.CErrorMessage, "failed: message is null")
			return
		}

		friendPubKey := asymmetric.LoadRSAPubKey(r.Header.Get(hls_settings.CHeaderPubKey))
		if friendPubKey == nil {
			panic("public key is null (invalid data from HLS)!")
		}

		myPubKey, err := client.PubKey()
		if err != nil {
			response(w, hls_settings.CErrorPubKey, "failed: message is null")
			return
		}

		rel := database.NewRelation(myPubKey, friendPubKey)
		dbMsg := database.NewMessage(true, msg)

		if err := db.Push(rel, dbMsg); err != nil {
			response(w, hls_settings.CErrorWrite, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:   friendPubKey.Address().String(),
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
