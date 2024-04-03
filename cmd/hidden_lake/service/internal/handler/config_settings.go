package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
)

func HandleConfigSettingsAPI(pWrapper config.IWrapper, pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

			sett := pWrapper.GetConfig().GetSettings()
			_ = api.Response(pW, http.StatusOK, pkg_config.SConfigSettings{
				SConfigSettings: config.SConfigSettings{
					FMessageSizeBytes:   sett.GetMessageSizeBytes(),
					FWorkSizeBits:       sett.GetWorkSizeBits(),
					FQueuePeriodMS:      sett.GetQueuePeriodMS(),
					FQueueRandPeriodMS:  sett.GetQueueRandPeriodMS(),
					FKeySizeBits:        sett.GetKeySizeBits(),
					FLimitVoidSizeBytes: sett.GetLimitVoidSizeBytes(),
					FNetworkKey:         sett.GetNetworkKey(),
					FF2FDisabled:        sett.GetF2FDisabled(),
				},
				FLimitMessageSizeBytes: pNode.GetMessageQueue().GetClient().GetMessageLimit(),
			})
			return

		case http.MethodPost:
			networkKeyBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: read network key bytes")
				return
			}

			networkKey := string(networkKeyBytes)
			if err := pWrapper.GetEditor().UpdateNetworkKey(networkKey); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_key"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: update network key")
				return
			}

			pNode.GetNetworkNode().SetVSettings(
				conn.NewVSettings(&conn.SVSettings{FNetworkKey: networkKey}),
			)
			pNode.GetMessageQueue().SetVSettings(
				queue.NewVSettings(&queue.SVSettings{FNetworkKey: networkKey}),
			)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: set network key")
			return
		}
	}
}
