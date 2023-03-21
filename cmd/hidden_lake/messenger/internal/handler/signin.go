package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

func SignInPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/sign/in" {
			NotFoundPage(pState)(pW, pR)
			return
		}

		if pState.IsActive() {
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

			hashLogin := hashing.NewSHA256Hasher([]byte(login))
			hashPassword := hashing.NewSHA256Hasher([]byte(password))

			hashLP := hashing.NewSHA256Hasher(bytes.Join(
				[][]byte{hashLogin.ToBytes(), hashPassword.ToBytes()},
				[]byte{},
			)).ToBytes()

			if err := pState.UpdateState(hashLP); err != nil {
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
		t.Execute(pW, pState.GetTemplate())
	}
}
