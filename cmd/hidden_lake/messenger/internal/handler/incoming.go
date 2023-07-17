package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderOffResponse, "true")

		if pR.Method != http.MethodPost {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if !pStateManager.StateIsActive() {
			api.Response(pW, http.StatusUnauthorized, "failed: client unauthorized")
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		rawMsg := strings.TrimSpace(string(msgBytes))
		if len(rawMsg) == 0 {
			api.Response(pW, http.StatusTeapot, "failed: message is null")
			return
		}

		if utils.HasNotWritableCharacters(rawMsg) {
			api.Response(pW, http.StatusUnsupportedMediaType, "failed: message has not writable characters")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is null (invalid data from HLS)!")
		}

		msgHash := pR.Header.Get(hls_settings.CHeaderMessageHash)
		if msgHash == "" {
			panic("message hash is null (invalid data from HLS)!")
		}

		myPubKey, err := pStateManager.GetClient().Service().GetPubKey()
		if err != nil {
			api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		db := pStateManager.GetWrapperDB().Get()
		if db == nil {
			api.Response(pW, http.StatusForbidden, "failed: database closed")
			return
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, utils.ReplaceTextToEmoji(rawMsg), encoding.HexDecode(msgHash))

		if err := db.Push(rel, dbMsg); err != nil {
			api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:   fPubKey.GetAddress().ToString(),
			FMessage:   dbMsg.GetMessage(),
			FTimestamp: dbMsg.GetTimestamp(),
		})
		api.Response(pW, http.StatusOK, pkg_settings.CTitlePattern)
	}
}
