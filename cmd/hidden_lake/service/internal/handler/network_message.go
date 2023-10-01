package handler

import (
	"io"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleNetworkMessageAPI(pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: read message bytes")
			return
		}

		msg := message.LoadMessage(
			pNode.GetMessageQueue().GetClient().GetSettings(),
			string(msgBytes),
		)
		if msg == nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_message"))
			api.Response(pW, http.StatusBadRequest, "failed: decode message")
			return
		}

		pNode.HandleMessage(msg)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, "success: handle message")
	}
}
