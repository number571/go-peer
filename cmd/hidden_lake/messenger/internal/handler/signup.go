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
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func SignUpPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sign/up" {
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

			passwordRepeat := strings.TrimSpace(r.FormValue("password_repeat"))
			if passwordRepeat == "" {
				fmt.Fprint(w, "error: password_repeat is null")
				return
			}

			if password != passwordRepeat {
				fmt.Fprint(w, "error: passwords not equals")
				return
			}

			hashLogin := hashing.NewSHA256Hasher([]byte(login))
			hashPassword := hashing.NewSHA256Hasher([]byte(password))

			hashLP := hashing.NewSHA256Hasher(bytes.Join(
				[][]byte{hashLogin.Bytes(), hashPassword.Bytes()},
				[]byte{},
			)).Bytes()

			var privKey asymmetric.IPrivKey
			privateKey := strings.TrimSpace(r.FormValue("private_key"))

			switch privateKey {
			case "":
				privKey = asymmetric.NewRSAPrivKey(hls_settings.CAKeySize)
			default:
				privKey = asymmetric.LoadRSAPrivKey(privateKey)
			}

			if privKey == nil {
				fmt.Fprint(w, "error: incorrect private key")
				return
			}

			if err := s.CreateState(hashLP, privKey); err != nil {
				fmt.Fprint(w, "error: create account")
				return
			}

			http.Redirect(w, r, "/sign/in", http.StatusFound)
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
		t.Execute(w, s.GetTemplate())
	}
}
