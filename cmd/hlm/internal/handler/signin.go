package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/entropy"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/storage"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
	"github.com/number571/go-peer/cmd/hlm/web"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
)

func SignInPage(wDB database.IWrapperDB, client hls_client.IClient, stg storage.IKeyValueStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/sign/in" {
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

			hashLogin := hashing.NewSHA256Hasher([]byte(login))
			hashPassword := hashing.NewSHA256Hasher([]byte(password))

			hashLP := hashing.NewSHA256Hasher(bytes.Join(
				[][]byte{hashLogin.Bytes(), hashPassword.Bytes()},
				[]byte{},
			)).Bytes()

			privKeyBytes, err := stg.Get(hashLP)
			if err != nil {
				fmt.Fprint(w, "error: account does not exist")
				return
			}

			privKey := asymmetric.LoadRSAPrivKey(privKeyBytes)
			if privKey == nil {
				fmt.Fprint(w, "error: private key is null")
				return
			}

			if err := client.PrivKey(privKey); err != nil {
				fmt.Fprint(w, "error: update private key")
				return
			}

			err = wDB.Update(database.NewKeyValueDB(
				hlm_settings.CPathDB,
				entropy.NewEntropy(hlm_settings.CWorkForKeys).
					Raise([]byte(password), []byte(login)),
			))
			if err != nil {
				fmt.Fprint(w, "error: update database")
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
			fmt.Println(err)
			panic("can't load hmtl files")
		}
		t.Execute(w, newTemplateData(db))
	}
}
