package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/stringtools"
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
			connects := stringtools.UniqAppendToSlice(
				pWrapper.GetConfig().GetConnections(),
				connect,
			)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: update connections")
				return
			}
			_ = pNode.GetNetworkNode().AddConnection(connect) // connection may be refused (closed)
			api.Response(pW, http.StatusOK, "success: update connections")
		case http.MethodDelete:
			connects := stringtools.DeleteFromSlice(pWrapper.GetConfig().GetConnections(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: delete connection")
				return
			}
			if err := pNode.GetNetworkNode().DelConnection(connect); err != nil {
				api.Response(pW, http.StatusInternalServerError, "failed: del connection")
				return
			}
			api.Response(pW, http.StatusOK, "success: delete connection")
		}
	}
}
