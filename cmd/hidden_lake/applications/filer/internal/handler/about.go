package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/web"
	"github.com/number571/go-peer/pkg/logger"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
)

func AboutPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.URL.Path != "/about" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
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
		t.Execute(pW, getTemplate(pCfg))
	}
}
