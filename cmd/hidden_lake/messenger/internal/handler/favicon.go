package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
)

func FaviconPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/favicon.ico" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		http.Redirect(pW, pR, "/static/img/favicon.ico", http.StatusFound)
	}
}
