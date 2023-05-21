package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigConnectsAPI(pWrapper config.IWrapper, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		switch pR.Method {
		case http.MethodGet:
			api.Response(pW, http.StatusOK, strings.Join(pWrapper.GetConfig().GetConnections(), ","))
			return
		case http.MethodPost, http.MethodDelete:
			// next
		default:
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		connectBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			api.Response(pW, http.StatusConflict, "failed: read connect bytes")
			return
		}

		connect := string(connectBytes)
		switch pR.Method {
		case http.MethodPost:
			connects := append(pWrapper.GetConfig().GetConnections(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: update connections")
				return
			}
			pNode.GetNetworkNode().AddConnect(connect)
			api.Response(pW, http.StatusOK, "success: update connections")
		case http.MethodDelete:
			connects := deleteConnect(pWrapper.GetConfig(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: delete connection")
				return
			}
			pNode.GetNetworkNode().DelConnect(connect)
			api.Response(pW, http.StatusOK, "success: delete connection")
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
