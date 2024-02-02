package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"

	internal_anon_logger "github.com/number571/go-peer/internal/logger/anon"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func HandleServiceTCP(pCfgW config.IWrapper, pLogger logger.ILogger) anonymity.IHandlerF {
	httpClient := &http.Client{Timeout: hls_settings.CFetchTimeout}

	return func(pCtx context.Context, pNode anonymity.INode, sender asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
		logBuilder := anon_logger.NewLogBuilder(hls_settings.CServiceName)

		// enrich logger
		logBuilder.
			WithSize(len(reqBytes)).
			WithPubKey(sender)

		cfg := pCfgW.GetConfig()
		friends := cfg.GetFriends()

		// append public key to list of friends if f2f option is disabled
		if cfg.GetSettings().GetF2FDisabled() && !inFriendsList(friends, sender) {
			// update config state with new friend
			friends[sender.GetHasher().ToString()] = sender
			if err := pCfgW.GetEditor().UpdateFriends(friends); err != nil {
				pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogBaseAppendNewFriend))
				return nil, err
			}
			// update list of friends and continue read request
			pNode.GetListPubKeys().AddPubKey(sender)
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseAppendNewFriend))
		}

		// load request from message's body
		loadReq, err := request.LoadRequest(reqBytes)
		if err != nil {
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroLoadRequestType))
			return nil, err
		}

		// get service's address by hostname
		service, ok := cfg.GetService(loadReq.GetHost())
		if !ok {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnUndefinedService))
			return nil, fmt.Errorf("failed: address undefined")
		}

		// generate new request to serivce
		pushReq, err := http.NewRequestWithContext(
			pCtx,
			loadReq.GetMethod(),
			fmt.Sprintf("http://%s%s", service.GetHost(), loadReq.GetPath()),
			bytes.NewReader(loadReq.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroProxyRequestType))
			return nil, err
		}

		// append headers from request & set service headers
		for key, val := range loadReq.GetHead() {
			pushReq.Header.Set(key, val)
		}
		pushReq.Header.Set(hls_settings.CHeaderPublicKey, sender.ToString())

		// send request to service
		// and receive response from service
		resp, err := httpClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestToService))
			return nil, err
		}
		defer resp.Body.Close()

		// get response mode: on/off
		respMode := resp.Header.Get(hls_settings.CHeaderResponseMode)
		switch respMode {
		case "", hls_settings.CHeaderResponseModeON:
			// send response to the client
			pLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoResponseFromService))
			return response.NewResponse(resp.StatusCode).
					WithHead(getResponseHead(resp)).
					WithBody(getResponseBody(resp)).
					ToBytes(),
				nil
		case hls_settings.CHeaderResponseModeOFF:
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

func inFriendsList(pFriends map[string]asymmetric.IPubKey, pPubKey asymmetric.IPubKey) bool {
	pubKey, ok := pFriends[pPubKey.GetHasher().ToString()]
	if !ok || !bytes.Equal(pubKey.ToBytes(), pPubKey.ToBytes()) {
		// the same keys, but different values
		return false
	}
	return true
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
