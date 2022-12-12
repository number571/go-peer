package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
)

type sSettings struct {
	FPublicKey   string
	FConnections []string
}

func SettingsPage(client hls_client.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/settings" {
			NotFoundPage(w, r)
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

		pubKey, err := client.PubKey()
		if err != nil {
			fmt.Fprint(w, "error: read connections")
			return
		}
		result.FPublicKey = pubKey.String()

		res, err := client.GetConnections()
		if err != nil {
			fmt.Fprint(w, "error: read connections")
			return
		}
		result.FConnections = res

		t, err := template.ParseFiles(
			hlm_settings.CPathTemplates+"index.html",
			hlm_settings.CPathTemplates+"settings.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, result)
	}
}
