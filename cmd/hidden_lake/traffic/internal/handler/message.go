package handler

import (
	"io"
	"net/http"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/internal/msgconv"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

func HandleMessageAPI(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pHTTPLogger, pAnonLogger logger.ILogger, pNode network.INode) http.HandlerFunc {
	tcpHandler := HandleServiceTCP(pCfg, pWrapperDB, pAnonLogger)

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlt_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()
		if database == nil {
			pHTTPLogger.PushErro(logBuilder.WithMessage("get_database"))
			api.Response(pW, http.StatusInternalServerError, "failed: get database")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			query := pR.URL.Query()
			msg, err := database.Load(query.Get("hash"))
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("load_hash"))
				api.Response(pW, http.StatusNotFound, "failed: load message")
				return
			}

			pHTTPLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, msg.ToString())
			return

		case http.MethodPost:
			msgStringAsBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			msgString := string(msgStringAsBytes)
			netMsg := net_message.NewMessage(
				pNode.GetSettings().GetConnSettings(),
				payload.NewPayload(
					hls_settings.CNetworkMask,
					msgconv.FromStringToBytes(msgString),
				),
			)

			if err := tcpHandler(pNode, nil, netMsg); err != nil {
				// internal logger
				api.Response(pW, http.StatusBadRequest, "failed: handle message")
				return
			}

			pHTTPLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: handle message")
			return
		}
	}
}
