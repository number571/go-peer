package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/pkg/errors"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
)

type sConnection struct {
	FAddress  string
	FIsBackup bool
	FOnline   bool
}

type sSettings struct {
	*state.STemplateState
	FPublicKey   string
	FConnections []sConnection
}

func SettingsPage(pStateManager state.IStateManager, pEditor config.IEditor) http.HandlerFunc {
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
		case http.MethodPut:
			language := strings.TrimSpace(pR.FormValue("language"))
			res, err := utils.ToILanguage(language)
			if err != nil {
				fmt.Fprint(pW, "error: load unknown language")
				return
			}
			if err := pEditor.UpdateLanguage(res); err != nil {
				fmt.Fprint(pW, "error: update language")
				return
			}
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
			err := pStateManager.AddConnection(
				fmt.Sprintf("%s:%s", host, port),
				strings.TrimSpace(pR.FormValue("is_backup")) != "",
			)
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

		client := pStateManager.GetClient()

		result := new(sSettings)
		result.STemplateState = pStateManager.GetTemplate()

		pubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}
		result.FPublicKey = pubKey.ToString()

		allConns, err := pStateManager.GetConnections()
		if err != nil {
			fmt.Fprint(pW, errors.WrapError(err, "error: get connections"))
			return
		}

		// HLS connections
		connsWithOnline, err := getOnlineConnections(client, allConns)
		if err != nil {
			fmt.Fprint(pW, errors.WrapError(err, "error: get online connections"))
			return
		}

		result.FConnections = connsWithOnline
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

func getOnlineConnections(client hls_client.IClient, allConns []state.IConnection) ([]sConnection, error) {
	var connections []sConnection

	onlines, err := client.GetOnlines()
	if err != nil {
		return nil, fmt.Errorf("error: read online connections")
	}

	for _, c := range allConns {
		connections = append(
			connections,
			sConnection{
				FAddress:  c.GetAddress(),
				FIsBackup: c.IsBackup(),
				FOnline:   getOnline(onlines, c),
			},
		)
	}

	return connections, nil
}

func getOnline(onlines []string, c state.IConnection) bool {
	if c.IsBackup() {
		return false
	}
	for _, o := range onlines {
		if o == c.GetAddress() {
			return true
		}
	}
	return false
}
