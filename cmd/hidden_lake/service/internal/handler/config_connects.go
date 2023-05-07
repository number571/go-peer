package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigConnectsAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			api.Response(pW, pkg_settings.CErrorNone, strings.Join(pWrapper.GetConfig().GetConnections(), ","))
			return
		}

		connectBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			api.Response(pW, pkg_settings.CErrorRead, "failed: read connect bytes")
			return
		}

		connect := string(connectBytes)
		switch pR.Method {
		case http.MethodPost:
			connects := append(pWrapper.GetConfig().GetConnections(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, pkg_settings.CErrorAction, "failed: update connections")
				return
			}
			pNode.GetNetworkNode().AddConnect(connect)
			api.Response(pW, pkg_settings.CErrorNone, "success: update connections")
		case http.MethodDelete:
			connects := deleteConnect(pWrapper.GetConfig(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, pkg_settings.CErrorAction, "failed: delete connection")
				return
			}
			pNode.GetNetworkNode().DelConnect(connect)
			api.Response(pW, pkg_settings.CErrorNone, "success: delete connection")
		}
	}
}
func deleteConnect(pCfg config.IConfig, pConnect string) []string {
	connects := pCfg.GetConnections()
	result := make([]string, 0, len(connects))
	for _, conn := range connects {
		if conn == pConnect {
			continue
		}
		result = append(result, conn)
	}
	return result
}
