package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"

	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func HandleServiceTCP(pCfg config.IConfig, pLogger logger.ILogger) anonymity.IHandlerF {
	logger := anon_logger.NewLogger(pkg_settings.CServiceName)

	return func(_ anonymity.INode, sender asymmetric.IPubKey, msgHash, reqBytes []byte) ([]byte, error) {
		// load request from message's body
		loadReq, err := request.LoadRequest(reqBytes)
		if err != nil {
			pLogger.PushErro(logger.GetFmtLog(pkg_settings.CLogErroLoadRequestType, msgHash, 0, sender, nil))
			return nil, err
		}

		// get service's address by hostname
		address, ok := pCfg.GetService(loadReq.GetHost())
		if !ok {
			pLogger.PushWarn(logger.GetFmtLog(pkg_settings.CLogWarnUndefinedService, msgHash, 0, sender, nil))
			return nil, fmt.Errorf("failed: address undefined")
		}

		// generate new request to serivce
		pushReq, err := http.NewRequest(
			loadReq.GetMethod(),
			fmt.Sprintf("http://%s%s", address, loadReq.GetPath()),
			bytes.NewReader(loadReq.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logger.GetFmtLog(pkg_settings.CLogErroProxyRequestType, msgHash, 0, sender, nil))
			return nil, err
		}

		// append headers from request & set service headers
		for key, val := range loadReq.GetHead() {
			pushReq.Header.Set(key, val)
		}
		pushReq.Header.Set(pkg_settings.CHeaderPublicKey, sender.ToString())
		pushReq.Header.Set(pkg_settings.CHeaderMessageHash, encoding.HexEncode(msgHash))

		// send request to service
		// and receive response from service
		resp, err := http.DefaultClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logger.GetFmtLog(pkg_settings.CLogWarnRequestToService, msgHash, 0, sender, nil))
			return nil, err
		}
		defer resp.Body.Close()

		// the response is not required by the client side
		if resp.Header.Get(pkg_settings.CHeaderOffResponse) != "" {
			pLogger.PushInfo(logger.GetFmtLog(pkg_settings.CLogInfoOffResponseFromService, msgHash, 0, sender, nil))
			return nil, nil
		}

		// send result to client
		pLogger.PushInfo(logger.GetFmtLog(pkg_settings.CLogInfoResponseFromService, msgHash, 0, sender, nil))
		return response.NewResponse(resp.StatusCode).
				WithHead(getHead(resp)).
				WithBody(getBody(resp)).
				ToBytes(),
			nil
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
