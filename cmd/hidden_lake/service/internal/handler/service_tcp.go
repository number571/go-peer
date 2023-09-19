package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"

	"github.com/number571/go-peer/pkg/network/anonymity/logbuilder"
)

func HandleServiceTCP(pCfg config.IConfig, pLogger logger.ILogger) anonymity.IHandlerF {
	return func(_ anonymity.INode, sender asymmetric.IPubKey, reqBytes []byte) ([]byte, error) {
		logBuilder := logbuilder.NewLogBuilder(pkg_settings.CServiceName)

		// enrich logger
		logBuilder.WithSize(len(reqBytes))
		logBuilder.WithPubKey(sender)

		// load request from message's body
		loadReq, err := request.LoadRequest(reqBytes)
		if err != nil {
			pLogger.PushErro(logBuilder.Get(pkg_settings.CLogErroLoadRequestType))
			return nil, err
		}

		// get service's address by hostname
		address, ok := pCfg.GetService(loadReq.GetHost())
		if !ok {
			pLogger.PushWarn(logBuilder.Get(pkg_settings.CLogWarnUndefinedService))
			return nil, fmt.Errorf("failed: address undefined")
		}

		// generate new request to serivce
		pushReq, err := http.NewRequest(
			loadReq.GetMethod(),
			fmt.Sprintf("http://%s%s", address, loadReq.GetPath()),
			bytes.NewReader(loadReq.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.Get(pkg_settings.CLogErroProxyRequestType))
			return nil, err
		}

		// append headers from request & set service headers
		for key, val := range loadReq.GetHead() {
			pushReq.Header.Set(key, val)
		}
		pushReq.Header.Set(pkg_settings.CHeaderPublicKey, sender.ToString())

		// send request to service
		// and receive response from service
		resp, err := http.DefaultClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logBuilder.Get(pkg_settings.CLogWarnRequestToService))
			return nil, err
		}
		defer resp.Body.Close()

		// get response mode: on/off
		respMode := resp.Header.Get(pkg_settings.CHeaderResponseMode)
		switch respMode {
		case "", pkg_settings.CHeaderResponseModeON:
			// send response to the client
			pLogger.PushInfo(logBuilder.Get(pkg_settings.CLogInfoResponseFromService))
			return response.NewResponse(resp.StatusCode).
					WithHead(getHead(resp)).
					WithBody(getBody(resp)).
					ToBytes(),
				nil
		case pkg_settings.CHeaderResponseModeOFF:
			// response is not required by the client side
			pLogger.PushInfo(logBuilder.Get(pkg_settings.CLogBaseResponseModeFromService))
			return nil, nil
		default:
			// unknown response mode
			pLogger.PushErro(logBuilder.Get(pkg_settings.CLogBaseResponseModeFromService))
			return nil, fmt.Errorf("failed: got invalid value of header (response-mode)")
		}
	}
}

func getHead(pResp *http.Response) map[string]string {
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

func getBody(pResp *http.Response) []byte {
	data, err := io.ReadAll(pResp.Body)
	if err != nil {
		return nil
	}
	return data
}
