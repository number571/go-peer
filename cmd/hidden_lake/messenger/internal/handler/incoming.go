package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlm_settings.CServiceName, pR)

		pW.Header().Set(hls_settings.CHeaderOffResponse, "true")

		if pR.Method != http.MethodPost {
			pLogger.PushInfo(httpLogger.Get(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if !pStateManager.StateIsActive() {
			pLogger.PushInfo(httpLogger.Get(http_logger.CLogRedirect))
			api.Response(pW, http.StatusUnauthorized, "failed: client unauthorized")
			return
		}

		rawMsgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushInfo(httpLogger.Get(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		if err := isValidMsgBytes(rawMsgBytes); err != nil {
			pLogger.PushWarn(httpLogger.Get("recv_message"))
			api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
			return
		}

		myPubKey, _, err := pStateManager.GetClient().GetPubKey()
		if err != nil || !pStateManager.IsMyPubKey(myPubKey) {
			pLogger.PushWarn(httpLogger.Get("get_public_key"))
			api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		db := pStateManager.GetWrapperDB().Get()
		if db == nil {
			pLogger.PushErro(httpLogger.Get("get_database"))
			api.Response(pW, http.StatusForbidden, "failed: database closed")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is null (invalid data from HLS)!")
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, doMessageProcessor(rawMsgBytes))

		if err := db.Push(rel, dbMsg); err != nil {
			pLogger.PushErro(httpLogger.Get("push_message"))
			api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:     fPubKey.GetAddress().ToString(),
			FMessageInfo: getMessageInfo(dbMsg.GetMessage(), dbMsg.GetTimestamp()),
		})

		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, hlm_settings.CTitlePattern)
	}
}

func doMessageProcessor(msgBytes []byte) []byte {
	if isText(msgBytes) {
		return []byte(utils.ReplaceTextToEmoji(string(msgBytes)))
	}
	return msgBytes
}

func isValidMsgBytes(rawMsgBytes []byte) error {
	switch {
	case isText(rawMsgBytes):
		strMsg := strings.TrimSpace(unwrapText(rawMsgBytes))
		if strMsg == "" {
			return errors.NewError("failed: message is null")
		}
		if utils.HasNotWritableCharacters(strMsg) {
			return errors.NewError("failed: message has not writable characters")
		}
		return nil
	case isFile(rawMsgBytes):
		filename, msgBytes := unwrapFile(rawMsgBytes)
		if filename == "" || len(msgBytes) == 0 {
			return errors.NewError("failed: unwrap file")
		}
		return nil
	default:
		return errors.NewError("failed: unknown message type")
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
