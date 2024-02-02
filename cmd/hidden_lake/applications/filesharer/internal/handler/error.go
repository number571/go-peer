package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/web"
	"github.com/number571/go-peer/pkg/logger"

	http_logger "github.com/number571/go-peer/internal/logger/http"
)

type sError struct {
	*sTemplate
	FTitle   string
	FMessage string
}

func ErrorPage(pLogger logger.ILogger, pCfg config.IConfig, pTitle, pMessage string) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		pW.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"error.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushWarn(logBuilder.WithMessage(pTitle))
		_ = t.Execute(pW, &sError{
			sTemplate: getTemplate(pCfg),
			FTitle:    pTitle,
			FMessage:  pMessage,
		})
	}
}

func NotFoundPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return ErrorPage(pLogger, pCfg, "404_page", "page not found")
}
