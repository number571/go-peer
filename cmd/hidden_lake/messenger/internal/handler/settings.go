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

func SettingsPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/settings" {
			NotFoundPage(pState)(pW, pR)
			return
		}

		if !pState.IsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			host := strings.TrimSpace(pR.FormValue("host"))
			port := strings.TrimSpace(pR.FormValue("port"))
			if host == "" || port == "" {
				fmt.Fprint(pW, "error: host or port is null")
				return
			}
			if _, err := strconv.Atoi(port); err != nil {
				fmt.Fprint(pW, "error: port is not a number")
				return
			}
			err := pState.AddConnection(fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				fmt.Fprint(pW, "error: add connection")
				return
			}
		case http.MethodDelete:
			address := strings.TrimSpace(pR.FormValue("address"))
			if address == "" {
				fmt.Fprint(pW, "error: address is null")
				return
			}
			err := pState.DelConnection(address)
			if err != nil {
				fmt.Fprint(pW, "error: del connection")
				return
			}
		}

		client := pState.GetClient().Service()

		result := new(sSettings)
		result.STemplateState = pState.GetTemplate()

		pubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}
		result.FPublicKey = pubKey.ToString()

		res, err := client.GetConnections()
		if err != nil {
			fmt.Fprint(pW, "error: read connections")
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
		t.Execute(pW, result)
	}
}
