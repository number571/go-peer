package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

type sSettings struct {
	*state.STemplateState
	FPublicKey   string
	FConnections []string
}

func SettingsPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/settings" {
			NotFoundPage(s)(w, r)
			return
		}

		if !s.IsActive() {
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
			err := s.AddConnection(fmt.Sprintf("%s:%s", host, port))
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
			err := s.DelConnection(address)
			if err != nil {
				fmt.Fprint(w, "error: del connection")
				return
			}
		}

		client := s.GetClient().Service()

		result := new(sSettings)
		result.STemplateState = s.GetTemplate()

		pubKey, err := client.GetPubKey()
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
