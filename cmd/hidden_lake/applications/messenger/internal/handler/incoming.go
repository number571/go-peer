package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/chat_queue"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/database"
	hlm_utils "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/utils"

	hlm_request "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/request"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigHTTP(pLogger logger.ILogger, pCfg config.IConfig, pDB database.IKVDatabase, pQueuePusher queue_set.IQueuePusher) http.HandlerFunc {
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

		encMsgBytes, err := io.ReadAll(pR.Body)
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

		if err := shareMessage(pCfg, fPubKey, pQueuePusher, requestID, encMsgBytes); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("share_message"))
			// continue
		}

		rawMsgBytes, err := decryptMsgBytes(pCfg, fPubKey, senderID, encMsgBytes)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decrypt_message"))
			api.Response(pW, http.StatusConflict, "failed: decrypt message")
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

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, senderID, rawMsgBytes)

		if err := pDB.Push(rel, dbMsg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		gChatQueue.Push(&chat_queue.SMessage{
			FAddress: fPubKey.GetHasher().ToString(),
			FMessageInfo: getMessageInfo(
				true,
				dbMsg.GetSenderID(),
				dbMsg.GetMessage(),
				dbMsg.GetTimestamp(),
			),
		})

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, hlm_settings.CTitlePattern)
	}
}

func shareMessage(pCfg config.IConfig, pSender asymmetric.IPubKey, pQueuePusher queue_set.IQueuePusher, pRequestID string, pBody []byte) error {
	if !pCfg.GetShare() {
		return nil
	}

	// request already exist in the queue
	if ok := pQueuePusher.Push([]byte(pRequestID), []byte{}); !ok {
		return nil
	}

	hlsClient := getClient(pCfg)

	friends, err := hlsClient.GetFriends()
	if err != nil {
		return err
	}

	lenFriends := len(friends)
	chBufErr := make(chan error, lenFriends)

	wg := sync.WaitGroup{}
	wg.Add(lenFriends)

	for aliasName, pubKey := range friends {
		go func(aliasName string, pubKey asymmetric.IPubKey) {
			defer wg.Done()

			// do not send a request to the creator of the request
			if bytes.Equal(pubKey.ToBytes(), pSender.ToBytes()) {
				return
			}

			chBufErr <- hlsClient.BroadcastRequest(
				aliasName,
				hlm_request.NewPushRequest(pSender, pRequestID, pBody),
			)
		}(aliasName, pubKey)
	}

	wg.Wait()
	close(chBufErr)

	errList := make([]error, 0, lenFriends)
	for err := range chBufErr {
		errList = append(errList, err)
	}

	return utils.MergeErrors(errList...)
}

func decryptMsgBytes(pCfg config.IConfig, pPubKey asymmetric.IPubKey, senderID, encMsgBytes []byte) ([]byte, error) {
	aliasName := ""

	friends, err := getClient(pCfg).GetFriends()
	if err != nil {
		return nil, err
	}
	for k, v := range friends {
		if bytes.Equal(v.ToBytes(), pPubKey.ToBytes()) {
			aliasName = k
			break
		}
	}
	if aliasName == "" {
		return nil, errors.New("alias name not found")
	}

	// secret key can be = nil
	secretKey := pCfg.GetSecretKeys()[aliasName]

	authKey := keybuilder.NewKeyBuilder(1, []byte(hlm_settings.CAuthSalt)).Build(secretKey)
	cipherKey := keybuilder.NewKeyBuilder(1, []byte(hlm_settings.CCipherSalt)).Build(secretKey)

	decBytes := symmetric.NewAESCipher(cipherKey).DecryptBytes(encMsgBytes)
	if len(decBytes) < hashing.CSHA256Size {
		return nil, errors.New("failed decrypt bytes")
	}

	msgBytes := decBytes[hashing.CSHA256Size:]

	authBytes := bytes.Join([][]byte{senderID, msgBytes}, []byte{})
	newHash := hashing.NewHMACSHA256Hasher(authKey, authBytes).ToBytes()
	if !bytes.Equal(decBytes[:hashing.CSHA256Size], newHash) {
		return nil, errors.New("failed auth bytes")
	}

	return msgBytes, nil
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

func getMessageInfo(pEscape bool, pSenderID string, pRawMsgBytes []byte, pTimestamp string) hlm_utils.SMessageInfo {
	switch {
	case isText(pRawMsgBytes):
		msgData := unwrapText(pRawMsgBytes, pEscape)
		if msgData == "" {
			panic("message data = nil")
		}
		return hlm_utils.SMessageInfo{
			FSenderID:  pSenderID,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	case isFile(pRawMsgBytes):
		filename, msgData := unwrapFile(pRawMsgBytes, pEscape)
		if filename == "" || msgData == "" {
			panic("filename = nil OR message data = nil")
		}
		return hlm_utils.SMessageInfo{
			FSenderID:  pSenderID,
			FFileName:  filename,
			FMessage:   msgData,
			FTimestamp: pTimestamp,
		}
	default:
		panic("got unknown message type")
	}
}
