package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

func SignUpPage(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/sign/up" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if pStateManager.StateIsActive() {
			http.Redirect(pW, pR, "/about", http.StatusFound)
			return
		}

		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			login := strings.TrimSpace(pR.FormValue("login"))
			if login == "" {
				fmt.Fprint(pW, "error: login is null")
				return
			}

			password := strings.TrimSpace(pR.FormValue("password"))
			if password == "" {
				fmt.Fprint(pW, "error: password is null")
				return
			}

			passwordRepeat := strings.TrimSpace(pR.FormValue("password_repeat"))
			if passwordRepeat == "" {
				fmt.Fprint(pW, "error: password_repeat is null")
				return
			}

			if password != passwordRepeat {
				fmt.Fprint(pW, "error: passwords not equals")
				return
			}

			hashLogin := hashing.NewSHA256Hasher([]byte(login))
			hashPassword := hashing.NewSHA256Hasher([]byte(password))

			hashLP := hashing.NewSHA256Hasher(bytes.Join(
				[][]byte{hashLogin.ToBytes(), hashPassword.ToBytes()},
				[]byte{},
			)).ToBytes()

			var privKey asymmetric.IPrivKey
			privateKey := strings.TrimSpace(pR.FormValue("private_key"))

			switch privateKey {
			case "":
				privKey = asymmetric.NewRSAPrivKey(pStateManager.GetConfig().GetKeySize())
			default:
				privKey = asymmetric.LoadRSAPrivKey(privateKey)
			}

			if privKey == nil {
				fmt.Fprint(pW, "error: incorrect private key")
				return
			}

			if err := pStateManager.CreateState(hashLP, privKey); err != nil {
				fmt.Fprint(pW, "error: create account")
				return
			}

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
		t.Execute(pW, pStateManager.GetTemplate())
	}
}
