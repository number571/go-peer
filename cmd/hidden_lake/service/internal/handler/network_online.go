package handler

import (
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkOnlineAPI(pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
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

			api.Response(pW, http.StatusOK, strings.Join(inOnline, ","))
		case http.MethodDelete:
			connectBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				api.Response(pW, http.StatusConflict, "failed: read connect bytes")
				return
			}

			if err := pNode.GetNetworkNode().DelConnection(string(connectBytes)); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: delete online connection")
				return
			}

			api.Response(pW, http.StatusOK, "success: delete online connection")
		}
	}
}
