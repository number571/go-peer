package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkOnlineAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vConnect pkg_settings.SConnect

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			api.Response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		switch r.Method {
		case http.MethodGet:
			conns := node.GetNetworkNode().GetConnections()

			inOnline := make([]string, 0, len(conns))
			for addr := range conns {
				inOnline = append(inOnline, addr)
			}

			sort.SliceStable(inOnline, func(i, j int) bool {
				return inOnline[i] < inOnline[j]
			})

			api.Response(w, pkg_settings.CErrorNone, strings.Join(inOnline, ","))
		case http.MethodDelete:
			if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
				api.Response(w, pkg_settings.CErrorDecode, "failed: decode request")
				return
			}

			if err := node.GetNetworkNode().DelConnect(vConnect.FConnect); err != nil {
				api.Response(w, pkg_settings.CErrorNone, "failed: delete online connection")
				return
			}

			api.Response(w, pkg_settings.CErrorNone, "success: delete online connection")
		}
	}
}
