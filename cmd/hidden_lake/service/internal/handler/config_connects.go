package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigConnectsAPI(wrapper config.IWrapper, node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vConnect pkg_settings.SConnect

		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodDelete {
			api.Response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			api.Response(w, pkg_settings.CErrorNone, strings.Join(wrapper.GetConfig().GetConnections(), ","))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
			api.Response(w, pkg_settings.CErrorDecode, "failed: decode request")
			return
		}

		switch r.Method {
		case http.MethodPost:
			connects := append(wrapper.GetConfig().GetConnections(), vConnect.FConnect)
			if err := wrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(w, pkg_settings.CErrorAction, "failed: update connections")
				return
			}
			node.GetNetworkNode().AddConnect(vConnect.FConnect)
			api.Response(w, pkg_settings.CErrorNone, "success: update connections")
		case http.MethodDelete:
			connects := deleteConnect(wrapper.GetConfig(), vConnect.FConnect)
			if err := wrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(w, pkg_settings.CErrorAction, "failed: delete connection")
				return
			}
			node.GetNetworkNode().DelConnect(vConnect.FConnect)
			api.Response(w, pkg_settings.CErrorNone, "success: delete connection")
		}
	}
}
func deleteConnect(cfg config.IConfig, connect string) []string {
	connects := cfg.GetConnections()
	result := make([]string, 0, len(connects))
	for _, conn := range connects {
		if conn == connect {
			continue
		}
		result = append(result, conn)
	}
	return result
}
