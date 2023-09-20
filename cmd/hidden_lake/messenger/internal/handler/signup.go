package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/logger"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	pvalidator "github.com/wagslane/go-password-validator"
)

func SignUpPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/sign/up" {
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
			client := pStateManager.GetClient()
			sett, err := client.GetSettings()
			if err != nil {
				pLogger.PushWarn(httpLogger.Get("get_settings"))
				fmt.Fprint(pW, "error: get settings from HLS")
				return
			}

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

			passwordRepeat := strings.TrimSpace(pR.FormValue("password_repeat"))
			if passwordRepeat == "" {
				pLogger.PushWarn(httpLogger.Get("get_password_repeat"))
				fmt.Fprint(pW, "error: password_repeat is null")
				return
			}

			if password != passwordRepeat {
				pLogger.PushWarn(httpLogger.Get("incorrect_password"))
				fmt.Fprint(pW, "error: passwords not equals")
				return
			}

			if err := pvalidator.Validate(password, hlm_settings.CMinEntropy); err != nil {
				pLogger.PushWarn(httpLogger.Get("password_is_weak"))
				fmt.Fprint(pW, "error: password is weak")
				return
			}

			var privKey asymmetric.IPrivKey
			privateKey := strings.TrimSpace(pR.FormValue("private_key"))

			switch privateKey {
			case "":
				privKey = asymmetric.NewRSAPrivKey(sett.GetKeySizeBits())
			default:
				privKey = asymmetric.LoadRSAPrivKey(privateKey)
			}

			if privKey == nil {
				pLogger.PushWarn(httpLogger.Get("get_private_key"))
				fmt.Fprint(pW, "error: incorrect private key")
				return
			}

			hashLP := keybuilder.NewKeyBuilder(hlm_settings.CWorkForKeys, []byte(login)).Build(password)
			if err := pStateManager.CreateState(hashLP, privKey); err != nil {
				pLogger.PushWarn(httpLogger.Get("create_state"))
				fmt.Fprint(pW, "error: create account")
				return
			}

			pLogger.PushInfo(httpLogger.Get(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"signup.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		t.Execute(pW, pStateManager.GetTemplate())
	}
}
