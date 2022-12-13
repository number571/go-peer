package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	"github.com/number571/go-peer/cmd/hlm/web"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
)

type sSettings struct {
	*sTemplateData
	FPublicKey   string
	FConnections []string
}

func SettingsPage(wDB database.IWrapperDB, client hls_client.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/settings" {
			NotFoundPage(db)(w, r)
			return
		}

		if db == nil {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		r.ParseForm()

		switch r.FormValue("method") {
		case http.MethodPost:
			host := strings.TrimSpace(r.FormValue("host"))
			port := strings.TrimSpace(r.FormValue("port"))
			if host == "" || port == "" {
				fmt.Fprint(w, "error: host or port is null")
				return
			}
			if _, err := strconv.Atoi(port); err != nil {
				fmt.Fprint(w, "error: port is not a number")
				return
			}
			err := client.AddConnection(fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				fmt.Fprint(w, "error: add connection")
				return
			}
		case http.MethodDelete:
			address := strings.TrimSpace(r.FormValue("address"))
			if address == "" {
				fmt.Fprint(w, "error: address is null")
				return
			}
			err := client.DelConnection(address)
			if err != nil {
				fmt.Fprint(w, "error: del connection")
				return
			}
		}

		result := new(sSettings)
		result.sTemplateData = newTemplateData(db)

		pubKey, err := client.PubKey()
		if err != nil {
			fmt.Fprint(w, "error: read public key")
			return
		}
		result.FPublicKey = pubKey.String()

		res, err := client.GetConnections()
		if err != nil {
			fmt.Fprint(w, "error: read connections")
			return
		}
		result.FConnections = res

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"settings.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, result)
	}
}
