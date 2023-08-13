package handler

import (
	"io"
	"net/http"
	"sort"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkOnlineAPI(pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(httpLogger.Get(http_logger.CLogMethod))
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

			pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, inOnline)
		case http.MethodDelete:
			connectBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(httpLogger.Get(http_logger.CLogDecodeBody))
				api.Response(pW, http.StatusConflict, "failed: read connect bytes")
				return
			}

			if err := pNode.GetNetworkNode().DelConnection(string(connectBytes)); err != nil {
				pLogger.PushWarn(httpLogger.Get("del_connection"))
				api.Response(pW, http.StatusInternalServerError, "failed: delete online connection")
				return
			}

			pLogger.PushWarn(httpLogger.Get(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: delete online connection")
		}
	}
}
