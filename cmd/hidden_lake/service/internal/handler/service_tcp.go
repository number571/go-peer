package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/queue_set"

	internal_anon_logger "github.com/number571/go-peer/internal/logger/anon"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func HandleServiceTCP(pCfg config.IConfig, pLogger logger.ILogger) anonymity.IHandlerF {
	queueSet := queue_set.NewQueueSet(
		queue_set.NewSettings(&queue_set.SSettings{
			FCapacity: pkg_settings.CHandleRequestQueueSize,
		}),
	)

	return func(pCtx context.Context, pNode anonymity.INode, sender asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
		logBuilder := anon_logger.NewLogBuilder(pkg_settings.CServiceName)

		// enrich logger
		logBuilder.
			WithSize(len(reqBytes)).
			WithPubKey(sender)

		// load request from message's body
		loadReq, err := request.LoadRequest(reqBytes)
		if err != nil {
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroLoadRequestType))
			return nil, err
		}

		// get unique ID of request from the header
		loadReqHead := loadReq.GetHead()
		requestID, ok := loadReqHead[pkg_settings.CHeaderRequestId]
		if !ok || len(requestID) != pkg_settings.CHandleRequestIDSize {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedRequestID))
			return nil, err
		}

		// try store ID of request to the queue
		if ok := queueSet.Push([]byte(requestID), []byte{}); !ok {
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoRequestIDAlreadyExist))
			return nil, nil
		}

		// get service's address by hostname
		service, ok := pCfg.GetService(loadReq.GetHost())
		if !ok {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedService))
			return nil, fmt.Errorf("failed: address undefined")
		}

		// share request to all friends
		if service.GetShare() {
			friends := pCfg.GetFriends()

			wg := sync.WaitGroup{}
			wg.Add(len(friends))

			for _, pubKey := range friends {
				go func(pubKey asymmetric.IPubKey) {
					defer wg.Done()

					// do not send a request to the creator of the request
					if bytes.Equal(pubKey.ToBytes(), sender.ToBytes()) {
						return
					}

					// redirect request to another node
					_ = pNode.BroadcastPayload(
						pCtx,
						pubKey,
						adapters.NewPayload(pkg_settings.CServiceMask, reqBytes),
					)
				}(pubKey)
			}

			wg.Wait()
		}

		// host can be nil only if share=true
		// this validation rule in the config
		if service.GetHost() == "" {
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoOnlyShareRequest))
			return nil, nil
		}

		// generate new request to serivce
		pushReq, err := http.NewRequest(
			loadReq.GetMethod(),
			fmt.Sprintf("http://%s%s", service.GetHost(), loadReq.GetPath()),
			bytes.NewReader(loadReq.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroProxyRequestType))
			return nil, err
		}

		// append headers from request & set service headers
		for key, val := range loadReqHead {
			pushReq.Header.Set(key, val)
		}
		pushReq.Header.Set(pkg_settings.CHeaderPublicKey, sender.ToString())

		// send request to service
		// and receive response from service
		resp, err := http.DefaultClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestToService))
			return nil, err
		}
		defer resp.Body.Close()

		// get response mode: on/off
		respMode := resp.Header.Get(pkg_settings.CHeaderResponseMode)
		switch respMode {
		case "", pkg_settings.CHeaderResponseModeON:
			// send response to the client
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoResponseFromService))
			return response.NewResponse(resp.StatusCode).
					WithHead(getResponseHead(resp)).
					WithBody(getResponseBody(resp)).
					ToBytes(),
				nil
		case pkg_settings.CHeaderResponseModeOFF:
			// response is not required by the client side
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, nil
		default:
			// unknown response mode
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, fmt.Errorf("failed: got invalid value of header (response-mode)")
		}
	}
}

func getResponseHead(pResp *http.Response) map[string]string {
	headers := make(map[string]string)
	for k := range pResp.Header {
		switch strings.ToLower(k) {
		case "date", "content-length": // ignore deanonymizing & redundant headers
			continue
		default:
			headers[k] = pResp.Header.Get(k)
		}
	}
	return headers
}

func getResponseBody(pResp *http.Response) []byte {
	data, err := io.ReadAll(pResp.Body)
	if err != nil {
		return nil
	}
	return data
}
