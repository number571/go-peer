package handler

import (
	"io"
	"net/http"
	"sort"
	"strings"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkOnlineAPI(pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			conns := pNode.GetNetworkNode().GetConnections()

			inOnline := make([]string, 0, len(conns))
			for addr := range conns {
				inOnline = append(inOnline, addr)
			}

			sort.SliceStable(inOnline, func(i, j int) bool {
				return inOnline[i] < inOnline[j]
			})

			api.Response(pW, pkg_settings.CErrorNone, strings.Join(inOnline, ","))
		case http.MethodDelete:
			connectBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				api.Response(pW, pkg_settings.CErrorRead, "failed: read connect bytes")
				return
			}

			if err := pNode.GetNetworkNode().DelConnect(string(connectBytes)); err != nil {
				api.Response(pW, pkg_settings.CErrorNone, "failed: delete online connection")
				return
			}

			api.Response(pW, pkg_settings.CErrorNone, "success: delete online connection")
		}
	}
}
