package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			api.Response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if !s.IsActive() {
			api.Response(w, hls_settings.CErrorUnauth, "failed: client unauthorized")
			return
		}

		msgBytes, err := io.ReadAll(r.Body)
		if err != nil {
			api.Response(w, hls_settings.CErrorResponse, "failed: response message")
			return
		}

		msg := strings.TrimSpace(string(msgBytes))
		if len(msg) == 0 {
			api.Response(w, hls_settings.CErrorMessage, "failed: message is null")
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
			api.Response(w, hls_settings.CErrorPubKey, "failed: message is null")
			return
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, msg, encoding.HexDecode(msgHash))

		db := s.GetWrapperDB().Get()
		if err := db.Push(rel, dbMsg); err != nil {
			api.Response(w, hls_settings.CErrorWrite, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:   fPubKey.GetAddress().ToString(),
			FMessage:   dbMsg.GetMessage(),
			FTimestamp: dbMsg.GetTimestamp(),
		})
		api.Response(w, hls_settings.CErrorNone, settings.CTitlePattern)
	}
}
