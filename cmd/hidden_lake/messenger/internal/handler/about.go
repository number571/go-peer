package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
)

func AboutPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/about" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"about.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		t.Execute(pW, pStateManager.GetTemplate())
	}
}
