package handler

import (
	"net/http"

	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/config"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
)

func FaviconPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.URL.Path != "/favicon.ico" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		http.Redirect(pW, pR, "/static/img/favicon.ico", http.StatusFound)
	}
}
