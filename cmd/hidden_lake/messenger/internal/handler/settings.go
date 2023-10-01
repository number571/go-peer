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
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/logger"
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
	FNetworkKey  string
	FConnections []sConnection
}

func SettingsPage(pStateManager state.IStateManager, pEditor config.IEditor, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/settings" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		pR.ParseForm()

		client := pStateManager.GetClient()

		myPubKey, _, err := client.GetPubKey()
		if err != nil || !pStateManager.IsMyPubKey(myPubKey) {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			fmt.Fprint(pW, "error: read public key")
			return
		}

		switch pR.FormValue("method") {
		case http.MethodPatch:
			networkKey := strings.TrimSpace(pR.FormValue("network_key"))
			if err := client.SetNetworkKey(networkKey); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("set_network_key"))
				fmt.Fprint(pW, "error: update network key")
				return
			}
		case http.MethodPut:
			language := strings.TrimSpace(pR.FormValue("language"))
			res, err := utils.ToILanguage(language)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("to_language"))
				fmt.Fprint(pW, "error: load unknown language")
				return
			}
			if err := pEditor.UpdateLanguage(res); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_language"))
				fmt.Fprint(pW, "error: update language")
				return
			}
		case http.MethodPost:
			host := strings.TrimSpace(pR.FormValue("host"))
			port := strings.TrimSpace(pR.FormValue("port"))

			if host == "" || port == "" {
				pLogger.PushWarn(logBuilder.WithMessage("get_host_port"))
				fmt.Fprint(pW, "error: host or port is null")
				return
			}
			if _, err := strconv.Atoi(port); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("port_to_int"))
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
					pLogger.PushWarn(logBuilder.WithMessage("update_backup_connections"))
					fmt.Fprint(pW, errors.WrapError(err, "error: update backup connections"))
					return
				}
			case false:
				if err := client.AddConnection(connect); err != nil {
					pLogger.PushWarn(logBuilder.WithMessage("add_connection"))
					fmt.Fprint(pW, "error: add connection")
					return
				}
			}
		case http.MethodDelete:
			connect := strings.TrimSpace(pR.FormValue("address"))
			if connect == "" {
				pLogger.PushWarn(logBuilder.WithMessage("get_connection"))
				fmt.Fprint(pW, "error: connect is null")
				return
			}

			connects := stringtools.DeleteFromSlice(
				pStateManager.GetConfig().GetBackupConnections(),
				connect,
			)
			if err := pEditor.UpdateBackupConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("delete_backup_connection"))
				fmt.Fprint(pW, "failed: delete backup connection")
				return
			}

			if err := client.DelConnection(connect); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("get_ndelete_connectionetwork_key"))
				fmt.Fprint(pW, "error: del connection")
				return
			}
		}

		result := new(sSettings)
		result.STemplateState = pStateManager.GetTemplate()

		result.FPublicKey = myPubKey.ToString()

		networkKey, err := client.GetNetworkKey()
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_network_key"))
			fmt.Fprint(pW, "error: read network key")
			return
		}
		result.FNetworkKey = networkKey

		// append HLS connections to backup connections
		allConns, err := getAllConnections(pStateManager.GetConfig(), client)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_all_connections"))
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

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		t.Execute(pW, result)
	}
}

func getAllConnections(pConfig config.IConfig, pClient hls_client.IClient) ([]sConnection, error) {
	var connections []sConnection

	conns, err := pClient.GetConnections()
	if err != nil {
		return nil, fmt.Errorf("error: read connections")
	}

	onlines, err := pClient.GetOnlines()
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

	for _, c := range pConfig.GetBackupConnections() {
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
