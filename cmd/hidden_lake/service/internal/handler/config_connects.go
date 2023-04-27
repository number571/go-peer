package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigConnectsAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		var vConnect pkg_settings.SConnect

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			api.Response(pW, pkg_settings.CErrorNone, strings.Join(pWrapper.GetConfig().GetConnections(), ","))
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vConnect); err != nil {
			api.Response(pW, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			connects := append(pWrapper.GetConfig().GetConnections(), vConnect.FConnect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, pkg_settings.CErrorAction, "failed: update connections")
				return
			}
			pNode.GetNetworkNode().AddConnect(vConnect.FConnect)
			api.Response(pW, pkg_settings.CErrorNone, "success: update connections")
		case http.MethodDelete:
			connects := deleteConnect(pWrapper.GetConfig(), vConnect.FConnect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, pkg_settings.CErrorAction, "failed: delete connection")
				return
			}
			pNode.GetNetworkNode().DelConnect(vConnect.FConnect)
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
