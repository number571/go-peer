package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/receiver"
	hlm_utils "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/utils"

	hlm_client "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/client"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pDB database.IKVDatabase,
	pMsgReceiver receiver.IMessageReceiver,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		senderPseudonym := pR.Header.Get(hlm_settings.CHeaderPseudonym)
		if !hlm_utils.PseudonymIsValid(senderPseudonym) {
			pLogger.PushWarn(logBuilder.WithMessage("get_sender_pseudonym"))
			api.Response(pW, http.StatusUnauthorized, "failed: get pseudonym from messenger")
			return
		}

		rawMsgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is nil (invalid data from HLS)!")
		}

		requestID := pR.Header.Get(hlm_settings.CHeaderRequestId)
		if len(requestID) != hlm_settings.CRequestIDSize {
			pLogger.PushWarn(logBuilder.WithMessage("request_id_size"))
			api.Response(pW, http.StatusNotAcceptable, "failed: request id size")
			return
		}

		// request already exist in the queue
		if ok, err := pDB.PushRequestID([]byte(requestID)); err != nil {
			if ok {
				pLogger.PushWarn(logBuilder.WithMessage("request_id_exist"))
				api.Response(pW, http.StatusLocked, "failed: request_id already exist")
				return
			}
			pLogger.PushErro(logBuilder.WithMessage("push_request_id"))
			api.Response(pW, http.StatusServiceUnavailable, "failed: push request_id to database")
			return
		}

		if pCfg.GetSettings().GetShareEnabled() {
			err := shareMessage(pCfg, fPubKey, requestID, senderPseudonym, rawMsgBytes)
			switch err {
			case nil:
				pLogger.PushInfo(logBuilder.WithMessage("share_message"))
			default:
				pLogger.PushWarn(logBuilder.WithMessage("share_message"))
			}
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

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, senderPseudonym, rawMsgBytes)

		if err := pDB.Push(rel, dbMsg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		pMsgReceiver.Send(&receiver.SMessage{
			FAddress: fPubKey.GetHasher().ToString(),
			FMessageInfo: getMessageInfo(
				true,
				dbMsg.GetPseudonym(),
				dbMsg.GetMessage(),
				dbMsg.GetTimestamp(),
			),
		})

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, hlm_settings.CTitlePattern)
	}
}

func shareMessage(
	pCfg config.IConfig,
	pSender asymmetric.IPubKey,
	pRequestID string,
	senderPseudonym string,
	pBody []byte,
) error {
	hlsClient := getClient(pCfg)

	hlmClient := hlm_client.NewClient(
		hlm_client.NewBuilder(),
		hlm_client.NewRequester(hlsClient),
	)

	friends, err := hlsClient.GetFriends()
	if err != nil {
		return err
	}

	lenFriends := len(friends)

	wg := sync.WaitGroup{}
	wg.Add(lenFriends)

	errList := make([]error, lenFriends)
	i := 0

	for aliasName, pubKey := range friends {
		go func(i int, aliasName string, pubKey asymmetric.IPubKey) {
			defer wg.Done()

			// do not send a request to the creator of the request
			if bytes.Equal(pubKey.ToBytes(), pSender.ToBytes()) {
				return
			}

			errList[i] = hlmClient.PushMessage(aliasName, senderPseudonym, pRequestID, pBody)
		}(i, aliasName, pubKey)
		i++
	}

	wg.Wait()
	return utils.MergeErrors(errList...)
}

func isValidMsgBytes(rawMsgBytes []byte) error {
	switch {
	case isText(rawMsgBytes):
		strMsg := strings.TrimSpace(unwrapText(rawMsgBytes, true))
		if strMsg == "" {
			return errors.New("failed: message is nil")
		}
		if hlm_utils.HasNotWritableCharacters(strMsg) {
			return errors.New("failed: message has not writable characters")
		}
		return nil
	case isFile(rawMsgBytes):
		filename, msgBytes := unwrapFile(rawMsgBytes, true)
		if filename == "" || len(msgBytes) == 0 {
			return errors.New("failed: unwrap file")
		}
		return nil
	default:
		return errors.New("failed: unknown message type")
	}
}

func getMessageInfo(pEscape bool, pPseudonym string, pRawMsgBytes []byte, pTimestamp string) hlm_utils.SMessageInfo {
	switch {
	case isText(pRawMsgBytes):
		msgData := unwrapText(pRawMsgBytes, pEscape)
		if msgData == "" {
			panic("message data = nil")
		}
		return hlm_utils.SMessageInfo{
			FPseudonym: pPseudonym,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	case isFile(pRawMsgBytes):
		filename, msgData := unwrapFile(pRawMsgBytes, pEscape)
		if filename == "" || msgData == "" {
			panic("filename = nil OR message data = nil")
		}
		return hlm_utils.SMessageInfo{
			FPseudonym: pPseudonym,
			FFileName:  filename,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	default:
		panic("got unknown message type")
	}
}
