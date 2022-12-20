package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"

	"github.com/number571/go-peer/cmd/hlm/internal/app/state"
	"github.com/number571/go-peer/cmd/hlm/web"
)

func SignInPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sign/in" {
			NotFoundPage(s)(w, r)
			return
		}

		if s.IsActive() {
			http.Redirect(w, r, "/about", http.StatusFound)
			return
		}

		r.ParseForm()

		switch r.FormValue("method") {
		case http.MethodPost:
			login := strings.TrimSpace(r.FormValue("login"))
			if login == "" {
				fmt.Fprint(w, "error: login is null")
				return
			}

			password := strings.TrimSpace(r.FormValue("password"))
			if password == "" {
				fmt.Fprint(w, "error: password is null")
				return
			}

			hashLogin := hashing.NewSHA256Hasher([]byte(login))
			hashPassword := hashing.NewSHA256Hasher([]byte(password))

			hashLP := hashing.NewSHA256Hasher(bytes.Join(
				[][]byte{hashLogin.Bytes(), hashPassword.Bytes()},
				[]byte{},
			)).Bytes()

			if err := s.UpdateState(hashLP); err != nil {
				fmt.Fprintf(w, "error: %s", err.Error())
				return
			}

			http.Redirect(w, r, "/about", http.StatusFound)
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
		t.Execute(w, s.GetTemplate())
	}
}
