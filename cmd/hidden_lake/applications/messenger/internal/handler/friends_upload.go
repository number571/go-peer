package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/web"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

type sUploadFile struct {
	*sTemplate
	FAliasName    string
	FMessageLimit uint64
}

func FriendsUploadPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/friends/upload" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "get_alias_name", "alias name is nil")(pW, pR)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"upload.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		msgLimit, err := utils.GetMessageLimit(getClient(pCfg))
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_message_size", "get message size (limit)")(pW, pR)
			return
		}

		res := &sUploadFile{
			sTemplate:     getTemplate(pCfg),
			FAliasName:    aliasName,
			FMessageLimit: msgLimit,
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = t.Execute(pW, res)
	}
}
