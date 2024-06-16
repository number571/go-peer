package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

func HandleConfigSettingsAPI(pCfg config.IConfig, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		sett := pCfg.GetSettings()
		_ = api.Response(pW, http.StatusOK, config.SConfigSettings{
			FMessageSizeBytes:   sett.GetMessageSizeBytes(),
			FWorkSizeBits:       sett.GetWorkSizeBits(),
			FKeySizeBits:        sett.GetKeySizeBits(),
			FLimitVoidSizeBytes: sett.GetLimitVoidSizeBytes(),
			FMessagesCapacity:   sett.GetMessagesCapacity(),
			FNetworkKey:         sett.GetNetworkKey(),
			FStorageEnabled:     sett.GetStorageEnabled(),
		})
	}
}
