package handler

import (
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"

	http_logger "github.com/number571/go-peer/internal/logger/http"
)

func HandleNetworkPubKeyAPI(pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, pNode.GetMessageQueue().GetClient().GetPubKey().ToString())
	}
}
