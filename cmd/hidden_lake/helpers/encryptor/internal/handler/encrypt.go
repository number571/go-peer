package handler

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/config"
	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func HandleMessageEncryptAPI(
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pClient client.IClient,
	pParallel uint64,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hle_settings.CServiceName, pR)

		var vContainer hle_settings.SContainer

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vContainer); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		pubKey := asymmetric.LoadRSAPubKey(vContainer.FPublicKey)
		if pubKey == nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_pubkey"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: decode public key")
			return
		}

		bodyData := encoding.HexDecode(vContainer.FHexData)
		if bodyData == nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_hex_data"))
			_ = api.Response(pW, http.StatusTeapot, "failed: decode hex data")
			return
		}

		msg, err := pClient.EncryptMessage(
			pubKey,
			payload.NewPayload64(vContainer.FPldHead, bodyData).ToBytes(),
		)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("encrypt_payload"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: encrypt payload")
			return
		}

		cfgSett := pConfig.GetSettings()
		netMsg := net_message.NewMessage(
			net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{
					FWorkSizeBits: cfgSett.GetWorkSizeBits(),
					FNetworkKey:   cfgSett.GetNetworkKey(),
				}),
				FParallel:             pParallel,
				FRandMessageSizeBytes: cfgSett.GetRandMessageSizeBytes(),
			}),
			payload.NewPayload32(hls_settings.CNetworkMask, msg),
		)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, netMsg.ToString())
	}
}
