package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/utils"

	internal_anon_logger "github.com/number571/go-peer/internal/logger/anon"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func HandleServiceTCP(pCfgW config.IWrapper) anonymity.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}

	return func(
		pCtx context.Context,
		pNode anonymity.INode,
		pSender asymmetric.IPubKey,
		pReqBytes []byte,
	) ([]byte, error) {
		logger := pNode.GetLogger()
		logBuilder := anon_logger.NewLogBuilder(hls_settings.CServiceName)

		// enrich logger
		logBuilder.
			WithSize(len(pReqBytes)).
			WithPubKey(pSender)

		cfg := pCfgW.GetConfig()
		friends := cfg.GetFriends()

		// append public key to list of friends if f2f option is disabled
		if cfg.GetSettings().GetF2FDisabled() && !inFriendsList(friends, pSender) {
			// update config state with new friend
			friends[pSender.GetHasher().ToString()] = pSender
			if err := pCfgW.GetEditor().UpdateFriends(friends); err != nil {
				logger.PushErro(logBuilder.WithType(internal_anon_logger.CLogBaseAppendNewFriend))
				return nil, utils.MergeErrors(ErrUpdateFriends, err)
			}
			// update list of friends and continue read request
			pNode.GetListPubKeys().AddPubKey(pSender)
			logger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseAppendNewFriend))
		}

		// load request from message's body
		loadReq, err := request.LoadRequest(pReqBytes)
		if err != nil {
			logger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroLoadRequestType))
			return nil, utils.MergeErrors(ErrLoadRequest, err)
		}

		// get service's address by hostname
		service, ok := cfg.GetService(loadReq.GetHost())
		if !ok {
			logger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedService))
			return nil, ErrUndefinedService
		}

		// generate new request to serivce
		pushReq, err := http.NewRequestWithContext(
			pCtx,
			loadReq.GetMethod(),
			fmt.Sprintf("http://%s%s", service.GetHost(), loadReq.GetPath()),
			bytes.NewReader(loadReq.GetBody()),
		)
		if err != nil {
			logger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroProxyRequestType))
			return nil, utils.MergeErrors(ErrBuildRequest, err)
		}

		// append headers from request & set service headers
		for key, val := range loadReq.GetHead() {
			pushReq.Header.Set(key, val)
		}
		pushReq.Header.Set(hls_settings.CHeaderPublicKey, pSender.ToString())

		// send request and receive response from service
		resp, err := httpClient.Do(pushReq)
		if err != nil {
			logger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestToService))
			return nil, utils.MergeErrors(ErrBadRequest, err)
		}
		defer resp.Body.Close()

		// get response mode: on/off
		respMode := resp.Header.Get(hls_settings.CHeaderResponseMode)
		switch respMode {
		case "", hls_settings.CHeaderResponseModeON:
			// send response to the client
			logger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoResponseFromService))
			return response.NewResponse(resp.StatusCode).
					WithHead(getResponseHead(resp)).
					WithBody(getResponseBody(resp)).
					ToBytes(),
				nil
		case hls_settings.CHeaderResponseModeOFF:
			// response is not required by the client side
			logger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, nil
		default:
			// unknown response mode
			logger.PushErro(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, ErrInvalidResponseMode
		}
	}
}

func inFriendsList(pFriends map[string]asymmetric.IPubKey, pPubKey asymmetric.IPubKey) bool {
	for _, pubKey := range pFriends {
		if bytes.Equal(pubKey.ToBytes(), pPubKey.ToBytes()) {
			return true
		}
	}
	return false
}

func getResponseHead(pResp *http.Response) map[string]string {
	headers := make(map[string]string)
	for k := range pResp.Header {
		switch strings.ToLower(k) {
		case "date", "content-length": // ignore deanonymizing & redundant headers
			continue
		case strings.ToLower(hls_settings.CHeaderResponseMode): // delete HLS headers
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
