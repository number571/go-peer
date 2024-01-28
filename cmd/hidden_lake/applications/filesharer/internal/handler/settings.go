package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/web"
	"github.com/number571/go-peer/internal/language"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
)

type sConnection struct {
	FAddress  string
	FIsBackup bool
	FOnline   bool
}

type sSettings struct {
	*sTemplate
	FPseudonym     string
	FNetworkKey    string
	FPublicKey     string
	FPublicKeyHash string
	FConnections   []sConnection
}

func SettingsPage(pLogger logger.ILogger, pWrapper config.IWrapper) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		cfg := pWrapper.GetConfig()
		cfgEditor := pWrapper.GetEditor()

		if pR.URL.Path != "/settings" {
			NotFoundPage(pLogger, cfg)(pW, pR)
			return
		}

		pR.ParseForm()

		client := getClient(cfg)
		myPubKey, err := client.GetPubKey()
		if err != nil {
			ErrorPage(pLogger, cfg, "get_public_key", "read public key")(pW, pR)
			return
		}

		switch pR.FormValue("method") {
		case http.MethodPatch:
			networkKey := strings.TrimSpace(pR.FormValue("network_key"))
			if err := client.SetNetworkKey(networkKey); err != nil {
				ErrorPage(pLogger, cfg, "set_network_key", "update network key")(pW, pR)
				return
			}
		case http.MethodPut:
			strLang := strings.TrimSpace(pR.FormValue("language"))
			ilang, err := language.ToILanguage(strLang)
			if err != nil {
				ErrorPage(pLogger, cfg, "to_language", "load unknown language")(pW, pR)
				return
			}
			if err := cfgEditor.UpdateLanguage(ilang); err != nil {
				ErrorPage(pLogger, cfg, "update_language", "update language")(pW, pR)
				return
			}
		case http.MethodPost:
			host := strings.TrimSpace(pR.FormValue("host"))
			port := strings.TrimSpace(pR.FormValue("port"))

			if host == "" || port == "" {
				ErrorPage(pLogger, cfg, "get_host_port", "host or port is nil")(pW, pR)
				return
			}
			if _, err := strconv.Atoi(port); err != nil {
				ErrorPage(pLogger, cfg, "port_to_int", "port is not a number")(pW, pR)
				return
			}

			connect := fmt.Sprintf("%s:%s", host, port)
			if err := client.AddConnection(connect); err != nil {
				ErrorPage(pLogger, cfg, "add_connection", "add connection")(pW, pR)
				return
			}
		case http.MethodDelete:
			connect := strings.TrimSpace(pR.FormValue("address"))
			if connect == "" {
				ErrorPage(pLogger, cfg, "get_connection", "connect is nil")(pW, pR)
				return
			}

			if err := client.DelConnection(connect); err != nil {
				ErrorPage(pLogger, cfg, "del_connection", "delete connection")(pW, pR)
				return
			}
		}

		result := new(sSettings)
		result.sTemplate = getTemplate(cfg)

		result.FPublicKey = myPubKey.ToString()
		result.FPublicKeyHash = myPubKey.GetHasher().ToString()

		gotSettings, err := client.GetSettings()
		if err != nil {
			ErrorPage(pLogger, cfg, "get_network_key", "read network key")(pW, pR)
			return
		}

		result.FNetworkKey = gotSettings.GetNetworkKey()

		// append HLS connections to backup connections
		allConns, err := getAllConnections(cfg, client)
		if err != nil {
			ErrorPage(pLogger, cfg, "get_all_connections", "get online connections")(pW, pR)
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
