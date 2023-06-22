package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
)

type sConnection struct {
	FAddress string
	FOnline  bool
}

type sSettings struct {
	*state.STemplateState
	FPublicKey   string
	FConnections []sConnection
}

func SettingsPage(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/settings" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
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
			err := pStateManager.AddConnection(fmt.Sprintf("%s:%s", host, port))
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
			err := pStateManager.DelConnection(address)
			if err != nil {
				fmt.Fprint(pW, "error: del connection")
				return
			}
		}

		client := pStateManager.GetClient().Service()

		result := new(sSettings)
		result.STemplateState = pStateManager.GetTemplate()

		pubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}
		result.FPublicKey = pubKey.ToString()

		conns, err := getConnections(client)
		if err != nil {
			fmt.Fprint(pW, err.Error())
			return
		}

		result.FConnections = conns

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

func getConnections(client hls_client.IClient) ([]sConnection, error) {
	var connections []sConnection

	conns, err := client.GetConnections()
	if err != nil {
		return nil, fmt.Errorf("error: read connections")
	}

	onlines, err := client.GetOnlines()
	if err != nil {
		return nil, fmt.Errorf("error: read online connections")
	}

	for _, c := range conns {
		connections = append(
			connections,
			sConnection{
				FAddress: c,
				FOnline:  getState(onlines, c),
			},
		)
	}

	return connections, nil
}

func getState(onlines []string, c string) bool {
	for _, o := range onlines {
		if o == c {
			return true
		}
	}
	return false
}
