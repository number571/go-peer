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
	"github.com/number571/go-peer/pkg/errors"

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

		rawMsgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		msgBytes, err := getMessageBytesToRecv(rawMsgBytes)
		if err != nil {
			api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
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
		dbMsg := database.NewMessage(true, doMessageProcessor(msgBytes), encoding.HexDecode(msgHash))
		if err := db.Push(rel, dbMsg); err != nil {
			api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:     fPubKey.GetAddress().ToString(),
			FMessageInfo: getMessageInfo(dbMsg.GetMessage(), dbMsg.GetTimestamp()),
		})
		api.Response(pW, http.StatusOK, pkg_settings.CTitlePattern)
	}
}

func doMessageProcessor(msgBytes []byte) []byte {
	if isText(msgBytes) {
		return []byte(utils.ReplaceTextToEmoji(string(msgBytes)))
	}
	return msgBytes
}

func getMessageBytesToRecv(rawMsgBytes []byte) ([]byte, error) {
	switch {
	case isText(rawMsgBytes):
		strMsg := strings.TrimSpace(unwrapText(rawMsgBytes))
		if strMsg == "" {
			return nil, errors.NewError("failed: message is null")
		}
		if utils.HasNotWritableCharacters(strMsg) {
			return nil, errors.NewError("failed: message has not writable characters")
		}
		return wrapText(strMsg), nil
	case isFile(rawMsgBytes):
		filename, msgBytes := unwrapFile(rawMsgBytes)
		if filename == "" || len(msgBytes) == 0 {
			return nil, errors.NewError("failed: unwrap file")
		}
		return rawMsgBytes, nil
	default:
		return nil, errors.NewError("failed: unknown message type")
	}
}

func getMessageInfo(pRawMsgBytes []byte, pTimestamp string) utils.SMessageInfo {
	switch {
	case isText(pRawMsgBytes):
		msgData := unwrapText(pRawMsgBytes)
		if msgData == "" {
			panic("message data = nil")
		}
		return utils.SMessageInfo{
			FFileName:  "",
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	case isFile(pRawMsgBytes):
		filename, msgData := unwrapFile(pRawMsgBytes)
		if filename == "" || msgData == "" {
			panic("filename = nil OR message data = nil")
		}
		return utils.SMessageInfo{
			FFileName:  filename,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	default:
		panic("got unknown message type")
	}
}
