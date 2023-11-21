package handler

import (
	"net/http"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
)

func HandleIndexAPI(pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hll_settings.CServiceName, pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		api.Response(pW, http.StatusOK, hll_settings.CTitlePattern)
	}
}
