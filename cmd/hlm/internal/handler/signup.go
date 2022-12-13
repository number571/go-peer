package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/storage"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
	hls_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
)

func SignUpPage(wDB database.IWrapperDB, stg storage.IKeyValueStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/sign/up" {
			NotFoundPage(db)(w, r)
			return
		}

		if db != nil {
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

			if _, err := stg.Get(hashLP); err == nil {
				fmt.Fprint(w, "error: account already exist")
				return
			}

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

			if err := stg.Set(hashLP, privKey.Bytes()); err != nil {
				fmt.Fprint(w, "error: create storage account")
				return
			}

			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		t, err := template.ParseFiles(
			hlm_settings.CPathTemplates+"index.html",
			hlm_settings.CPathTemplates+"signup.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, newTemplateData(db))
	}
}
