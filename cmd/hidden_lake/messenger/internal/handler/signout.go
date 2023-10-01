package handler

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

func SignOutPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/sign/out" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		if err := pStateManager.CloseState(); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("close_state"))
			fmt.Fprint(pW, "error: clean hls_client data")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		http.Redirect(pW, pR, "/about", http.StatusFound)
	}
}
