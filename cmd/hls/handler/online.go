package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func HandleOnlineAPI(node anonymity.INode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vConnect hls_settings.SConnect

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			response(w, hls_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		if r.Method == http.MethodGet {
			conns := node.Network().Connections()
			inOnline := make([]string, 0, len(conns))
			for addr := range conns {
				inOnline = append(inOnline, addr)
			}
			sort.SliceStable(inOnline, func(i, j int) bool {
				return inOnline[i] < inOnline[j]
			})
			response(w, hls_settings.CErrorNone, strings.Join(inOnline, ","))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&vConnect); err != nil {
			response(w, hls_settings.CErrorDecode, "failed: decode request")
			return
		}

		if err := node.Network().Disconnect(vConnect.FConnect); err != nil {
			response(w, hls_settings.CErrorNone, "failed: delete online connection")
			return
		}
		response(w, hls_settings.CErrorNone, "success: delete online connection")
	}
}
