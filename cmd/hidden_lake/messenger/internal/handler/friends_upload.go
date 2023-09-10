package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
)

type sUploadFile struct {
	*state.STemplateState
	FAliasName    string
	FMessageLimit uint64
}

func FriendsUploadPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/friends/upload" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			pLogger.PushInfo(httpLogger.Get(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			pLogger.PushWarn(httpLogger.Get("get_alias_name"))
			fmt.Fprint(pW, "alias name is null")
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

		msgLimit, err := getMessageLimit(pStateManager.GetClient())
		if err != nil {
			pLogger.PushWarn(httpLogger.Get("get_message_size"))
			fmt.Fprint(pW, "get message size (limit)")
			return
		}

		res := &sUploadFile{
			STemplateState: pStateManager.GetTemplate(),
			FAliasName:     aliasName,
			FMessageLimit:  msgLimit,
		}

		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		t.Execute(pW, res)
	}
}
