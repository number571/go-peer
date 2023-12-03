package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkKeyAPI(pWrapper config.IWrapper, pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, pNode.GetNetworkNode().GetSettings().GetConnSettings().GetNetworkKey())
			return

		case http.MethodPost:
			networkKeyBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				api.Response(pW, http.StatusConflict, "failed: read network key bytes")
				return
			}

			networkKey := string(networkKeyBytes)
			if err := pWrapper.GetEditor().UpdateNetworkKey(networkKey); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_key"))
				api.Response(pW, http.StatusInternalServerError, "failed: update network key")
				return
			}

			pNode.GetNetworkNode().GetSettings().GetConnSettings().SetNetworkKey(networkKey)
			pNode.GetMessageQueue().ClearQueue()

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: set network key")
			return
		}
	}
}
