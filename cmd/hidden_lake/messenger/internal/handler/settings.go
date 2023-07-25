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
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/stringtools"

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

			connect := fmt.Sprintf("%s:%s", host, port)
			isBackup := strings.TrimSpace(pR.FormValue("is_backup")) != ""

			switch isBackup {
			case true:
				connects := stringtools.UniqAppendToSlice(
					pStateManager.GetConfig().GetBackupConnections(),
					connect,
				)
				if err := pEditor.UpdateBackupConnections(connects); err != nil {
					fmt.Fprint(pW, errors.WrapError(err, "error: update backup connections"))
					return
				}
			case false:
				err := pStateManager.AddConnection(connect)
				if err != nil {
					fmt.Fprint(pW, "error: add connection")
					return
				}
			}
		case http.MethodDelete:
			connect := strings.TrimSpace(pR.FormValue("address"))
			if connect == "" {
				fmt.Fprint(pW, "error: connect is null")
				return
			}

			connects := stringtools.DeleteFromSlice(
				pStateManager.GetConfig().GetBackupConnections(),
				connect,
			)
			if err := pEditor.UpdateBackupConnections(connects); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: delete backup connection")
				return
			}

			err := pStateManager.DelConnection(connect)
			if err != nil {
				fmt.Fprint(pW, "error: del connection")
				return
			}
		}

		result := new(sSettings)
		result.STemplateState = pStateManager.GetTemplate()

		client := pStateManager.GetClient()
		pubKey, err := client.GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}
		result.FPublicKey = pubKey.ToString()

		// HLS connections
		cfg := pStateManager.GetConfig()
		allConns, err := getAllConnections(client, cfg.GetBackupConnections())
		if err != nil {
			fmt.Fprint(pW, errors.WrapError(err, "error: get online connections"))
			return
		}

		result.FConnections = allConns
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

func getAllConnections(client hls_client.IClient, backupConns []string) ([]sConnection, error) {
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
				FOnline:  getOnline(onlines, c),
			},
		)
	}

	for _, c := range backupConns {
		connections = append(
			connections,
			sConnection{
				FAddress:  c,
				FIsBackup: true,
			},
		)
	}

	return connections, nil
}

func getOnline(onlines []string, c string) bool {
	for _, o := range onlines {
		if o == c {
			return true
		}
	}
	return false
}
