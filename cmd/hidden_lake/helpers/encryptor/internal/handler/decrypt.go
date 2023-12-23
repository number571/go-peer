package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"
	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func HandleDecryptAPI(pConfig config.IConfig, pLogger logger.ILogger, pClient client.IClient) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hle_settings.CServiceName, pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgStringAsBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: read encrypted message")
			return
		}

		netMsg, err := net_message.LoadMessage(pConfig.GetSettings(), string(msgStringAsBytes))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_net_message"))
			api.Response(pW, http.StatusNotAcceptable, "failed: decode network message")
			return
		}

		netPld := netMsg.GetPayload()
		if netPld.GetHead() != hls_settings.CNetworkMask {
			pLogger.PushWarn(logBuilder.WithMessage("invalid_net_mask"))
			api.Response(pW, http.StatusUnsupportedMediaType, "failed: invalid network mask")
			return
		}

		msg, err := message.LoadMessage(pConfig.GetSettings(), netPld.GetBody())
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_message"))
			api.Response(pW, http.StatusTeapot, "failed: decode message")
			return
		}

		pubKey, pld, err := pClient.DecryptMessage(msg)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decrypt_message"))
			api.Response(pW, http.StatusBadRequest, "failed: decrypt message")
			return
		}

		if uint32(pld.GetHead()) != hls_settings.CServiceMask {
			pLogger.PushWarn(logBuilder.WithMessage("invalid_service_mask"))
			api.Response(pW, http.StatusFailedDependency, "failed: invalid service mask")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, hle_settings.SContainer{
			FPublicKey: pubKey.ToString(),
			FHexData:   encoding.HexEncode(pld.GetBody()),
		})
	}
}
