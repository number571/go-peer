package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"
	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func HandleMessageDecryptAPI(pConfig config.IConfig, pLogger logger.ILogger, pClient client.IClient) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hle_settings.CServiceName, pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgStringAsBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: read encrypted message")
			return
		}

		netMsg, err := net_message.LoadMessage(
			pConfig.GetSettings(),
			net_message.CWithoutTimestamp,
			string(msgStringAsBytes),
		)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_net_message"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: decode network message")
			return
		}

		netPld := netMsg.GetPayload()
		if netPld.GetHead() != hls_settings.CNetworkMask {
			pLogger.PushWarn(logBuilder.WithMessage("invalid_net_mask"))
			_ = api.Response(pW, http.StatusUnsupportedMediaType, "failed: invalid network mask")
			return
		}

		pubKey, decMsg, err := pClient.DecryptMessage(netPld.GetBody())
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decrypt_message"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: decrypt message")
			return
		}

		pld := payload.LoadPayload64(decMsg)
		if pld == nil {
			pLogger.PushWarn(logBuilder.WithMessage("load_payload"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: load payload")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, hle_settings.SContainer{
			FPublicKey: pubKey.ToString(),
			FPldHead:   pld.GetHead(),
			FHexData:   encoding.HexEncode(pld.GetBody()),
		})
	}
}
