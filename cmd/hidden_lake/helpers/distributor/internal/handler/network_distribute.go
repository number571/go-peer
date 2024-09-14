package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/internal/config"
	hld_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
)

func HandleNetworkDistributeAPI(
	pCtx context.Context,
	pCfg config.IConfig,
	pLogger logger.ILogger,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hld_settings.CServiceName, pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		if pR.Method != http.MethodPost {
			pLogger.PushErro(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		reqBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		// load request from message's body
		loadReq, err := request.LoadRequest(string(reqBytes))
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("load_request"))
			_ = api.Response(pW, http.StatusPreconditionFailed, "load request")
			return
		}

		// get service's address by hostname
		service, ok := pCfg.GetService(loadReq.GetHost())
		if !ok {
			pLogger.PushWarn(logBuilder.WithMessage("get_service"))
			_ = api.Response(pW, http.StatusNotFound, "get service")
			return
		}

		// generate new request to serivce
		pushReq, err := http.NewRequestWithContext(
			pCtx,
			loadReq.GetMethod(),
			fmt.Sprintf("http://%s%s", service, loadReq.GetPath()),
			bytes.NewReader(loadReq.GetBody()),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("new_request"))
			_ = api.Response(pW, http.StatusInternalServerError, "new request")
			return
		}

		// append headers from request
		for key, val := range loadReq.GetHead() {
			pushReq.Header.Set(key, val)
		}

		// send request and receive response from service
		httpClient := &http.Client{Timeout: time.Minute}
		httpResp, err := httpClient.Do(pushReq)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("do_request"))
			_ = api.Response(pW, http.StatusBadRequest, "do request")
			return
		}
		defer httpResp.Body.Close()

		resp := response.NewResponse(httpResp.StatusCode).
			WithHead(getResponseHead(httpResp)).
			WithBody(getResponseBody(httpResp))

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, resp.ToString())
	}
}

func getResponseHead(pResp *http.Response) map[string]string {
	headers := make(map[string]string, len(pResp.Header))
	for k := range pResp.Header {
		if _, ok := gIgnoreHeaders[k]; ok {
			continue
		}
		headers[k] = pResp.Header.Get(k)
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
