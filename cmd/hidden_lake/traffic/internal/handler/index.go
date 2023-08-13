package handler

import (
	"net/http"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

func HandleIndexAPI(pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlt_settings.CServiceName, pR)
		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, hlt_settings.CTitlePattern)
	}
}
