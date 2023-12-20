package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(pLogger logger.ILogger, pCfg config.IConfig, pDB database.IKVDatabase) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		senderID := encoding.HexDecode(pR.Header.Get(hlm_settings.CHeaderSenderId))
		if len(senderID) != hashing.CSHA256Size {
			pLogger.PushWarn(logBuilder.WithMessage("get_sender_id"))
			api.Response(pW, http.StatusUnauthorized, "failed: get sender id from messenger")
			return
		}

		rawMsgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		if err := isValidMsgBytes(rawMsgBytes); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("recv_message"))
			api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
			return
		}

		myPubKey, err := getClient(pCfg).GetPubKey()
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is null (invalid data from HLS)!")
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, senderID, doMessageProcessor(rawMsgBytes))

		if err := pDB.Push(rel, dbMsg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress:     fPubKey.GetAddress().ToString(),
			FMessageInfo: getMessageInfo(dbMsg.GetSenderID(), dbMsg.GetMessage(), dbMsg.GetTimestamp()),
		})

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
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
			return errors.New("failed: message is null")
		}
		if utils.HasNotWritableCharacters(strMsg) {
			return errors.New("failed: message has not writable characters")
		}
		return nil
	case isFile(rawMsgBytes):
		filename, msgBytes := unwrapFile(rawMsgBytes)
		if filename == "" || len(msgBytes) == 0 {
			return errors.New("failed: unwrap file")
		}
		return nil
	default:
		return errors.New("failed: unknown message type")
	}
}

func getMessageInfo(pSenderID string, pRawMsgBytes []byte, pTimestamp string) utils.SMessageInfo {
	cutSenderID := fmt.Sprintf("%s...%s", string(pSenderID[:8]), string(pSenderID[len(pSenderID)-8:]))
	switch {
	case isText(pRawMsgBytes):
		msgData := unwrapText(pRawMsgBytes)
		if msgData == "" {
			panic("message data = nil")
		}
		return utils.SMessageInfo{
			FSenderID:  cutSenderID,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	case isFile(pRawMsgBytes):
		filename, msgData := unwrapFile(pRawMsgBytes)
		if filename == "" || msgData == "" {
			panic("filename = nil OR message data = nil")
		}
		return utils.SMessageInfo{
			FSenderID:  cutSenderID,
			FFileName:  filename,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	default:
		panic("got unknown message type")
	}
}
