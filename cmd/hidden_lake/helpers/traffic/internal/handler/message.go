package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func HandleMessageAPI(pCtx context.Context, pCfg config.IConfig, pDBWrapper database.IDBWrapper, pHTTPLogger, pAnonLogger logger.ILogger, pNode network.INode) http.HandlerFunc {
	tcpHandler := HandleServiceTCP(pCfg, pDBWrapper, pAnonLogger)

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlt_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		database := pDBWrapper.Get()
		if database == nil {
			pHTTPLogger.PushErro(logBuilder.WithMessage("get_database"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: get database")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			query := pR.URL.Query()

			hash := encoding.HexDecode(query.Get("hash"))
			if hash == nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusNotFound, "failed: decode message")
				return
			}

			msg, err := database.Load(hash)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("load_hash"))
				_ = api.Response(pW, http.StatusNotFound, "failed: load message")
				return
			}

			pHTTPLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, msg.ToString())
			return

		case http.MethodPost:
			msgStringAsBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			netMsg, err := net_message.LoadMessage(
				net_message.NewSettings(&net_message.SSettings{
					FNetworkKey:   pNode.GetVSettings().GetNetworkKey(),
					FWorkSizeBits: pNode.GetSettings().GetConnSettings().GetWorkSizeBits(),
				}),
				string(msgStringAsBytes),
			)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("decode_message"))
				_ = api.Response(pW, http.StatusTeapot, "failed: decode message")
				return
			}

			if netMsg.GetPayload().GetHead() != hls_settings.CNetworkMask {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("network_mask"))
				_ = api.Response(pW, http.StatusLocked, "failed: network mask")
				return
			}

			if err := tcpHandler(pCtx, pNode, nil, netMsg); err != nil {
				// internal logger
				_ = api.Response(pW, http.StatusBadRequest, "failed: handle message")
				return
			}

			pHTTPLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: handle message")
			return
		}
	}
}
