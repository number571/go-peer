package handler

import (
	"net/http"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hld_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
)

func HandleIndexAPI(pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hld_settings.CServiceName, pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		_ = api.Response(pW, http.StatusOK, hld_settings.CServiceFullName)
	}
}
