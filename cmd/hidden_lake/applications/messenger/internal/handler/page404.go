package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/web"
	"github.com/number571/go-peer/pkg/logger"

	http_logger "github.com/number571/go-peer/internal/logger/http"
)

func NotFoundPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		pW.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"page404.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogNotFound))
		t.Execute(pW, getTemplate(pCfg))
	}
}
