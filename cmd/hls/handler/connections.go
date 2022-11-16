package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hls/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func HandleConnectionsAPI(wrapper config.IWrapper, node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vConnect hls_settings.SConnect

		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodDelete {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			response(w, hls_settings.CErrorNone, strings.Join(wrapper.Config().Connections(), ","))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
			response(w, hls_settings.CErrorDecode, "failed: decode request")
			return
		}

		switch r.Method {
		case http.MethodPost:
			connects := append(wrapper.Config().Connections(), vConnect.FConnect)
			if err := wrapper.Editor().UpdateConnections(connects); err != nil {
				response(w, hls_settings.CErrorAction, "failed: update connections")
				return
			}
			node.Network().Connect(vConnect.FConnect)
			response(w, hls_settings.CErrorNone, "success: update connections")
		case http.MethodDelete:
			connects := deleteConnect(wrapper.Config(), vConnect.FConnect)
			if err := wrapper.Editor().UpdateConnections(connects); err != nil {
				response(w, hls_settings.CErrorAction, "failed: delete connection")
				return
			}
			node.Network().Disconnect(vConnect.FConnect)
			response(w, hls_settings.CErrorNone, "success: delete connection")
		}
	}
}
func deleteConnect(cfg config.IConfig, connect string) []string {
	connects := cfg.Connections()
	result := make([]string, 0, len(connects))
	for _, conn := range connects {
		if conn == connect {
			continue
		}
		result = append(result, conn)
	}
	return result
}
