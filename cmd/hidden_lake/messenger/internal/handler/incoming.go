package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if !s.IsActive() {
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

		fPubKey := asymmetric.LoadRSAPubKey(r.Header.Get(hls_settings.CHeaderPubKey))
		if fPubKey == nil {
			panic("public key is null (invalid data from HLS)!")
		}

		msgHash := r.Header.Get(hls_settings.CHeaderMsgHash)
		if msgHash == "" {
			panic("message hash is null (invalid data from HLS)!")
		}

		myPubKey, err := s.GetClient().Service().GetPubKey()
		if err != nil {
			response(w, hls_settings.CErrorPubKey, "failed: message is null")
			return
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, msg, encoding.HexDecode(msgHash))

		db := s.GetWrapperDB().Get()
		if err := db.Push(rel, dbMsg); err != nil {
			response(w, hls_settings.CErrorWrite, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:   fPubKey.Address().String(),
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
