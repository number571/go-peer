package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/msgbroker"
	hlm_utils "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pDB database.IKVDatabase,
	pBroker msgbroker.IMessageBroker,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		rawMsgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is nil (invalid data from HLS)!")
		}

		if err := isValidMsgBytes(rawMsgBytes); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("recv_message"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
			return
		}

		myPubKey, err := getHLSClient(pCfg).GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, rawMsgBytes)

		if err := pDB.Push(rel, dbMsg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		pBroker.Produce(
			fPubKey.GetHasher().ToString(),
			getMessage(
				true,
				dbMsg.GetMessage(),
				dbMsg.GetTimestamp(),
			),
		)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, hlm_settings.CServiceFullName)
	}
}

func isValidMsgBytes(rawMsgBytes []byte) error {
	switch {
	case isText(rawMsgBytes):
		strMsg := strings.TrimSpace(unwrapText(rawMsgBytes, true))
		if strMsg == "" {
			return ErrMessageNull
		}
		if hlm_utils.HasNotWritableCharacters(strMsg) {
			return ErrHasNotWritableChars
		}
		return nil
	case isFile(rawMsgBytes):
		filename, msgBytes := unwrapFile(rawMsgBytes, true)
		if filename == "" || len(msgBytes) == 0 {
			return ErrUnwrapFile
		}
		return nil
	default:
		return ErrUnknownMessageType
	}
}

func getMessage(pEscape bool, pRawMsgBytes []byte, pTimestamp string) hlm_utils.SMessage {
	switch {
	case isText(pRawMsgBytes):
		msgData := unwrapText(pRawMsgBytes, pEscape)
		if msgData == "" {
			panic("message data = nil")
		}
		return hlm_utils.SMessage{
			FTimestamp: pTimestamp,
			FMainData:  msgData,
		}
	case isFile(pRawMsgBytes):
		filename, msgData := unwrapFile(pRawMsgBytes, pEscape)
		if filename == "" || msgData == "" {
			panic("filename = nil OR message data = nil")
		}
		return hlm_utils.SMessage{
			FFileName:  filename,
			FTimestamp: pTimestamp,
			FMainData:  msgData,
		}
	default:
		panic("got unknown message type")
	}
}
