package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

func SignInPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/sign/in" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		if pStateManager.StateIsActive() {
			pLogger.PushInfo(httpLogger.Get(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/about", http.StatusFound)
			return
		}

		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			login := strings.TrimSpace(pR.FormValue("login"))
			if login == "" {
				pLogger.PushWarn(httpLogger.Get("get_login"))
				fmt.Fprint(pW, "error: login is null")
				return
			}

			password := strings.TrimSpace(pR.FormValue("password"))
			if password == "" {
				pLogger.PushWarn(httpLogger.Get("get_password"))
				fmt.Fprint(pW, "error: password is null")
				return
			}

			hashLogin := hashing.NewSHA256Hasher([]byte(login))
			hashPassword := hashing.NewSHA256Hasher([]byte(password))

			hashLP := hashing.NewSHA256Hasher(bytes.Join(
				[][]byte{hashLogin.ToBytes(), hashPassword.ToBytes()},
				[]byte{},
			)).ToBytes()

			if err := pStateManager.OpenState(hashLP); err != nil {
				pLogger.PushWarn(httpLogger.Get("open_state"))
				fmt.Fprintf(pW, "error: %s", err.Error())
				return
			}

			http.Redirect(pW, pR, "/about", http.StatusFound)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"signin.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		t.Execute(pW, pStateManager.GetTemplate())
	}
}
