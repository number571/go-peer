package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
)

func HandleConfigSettingsAPI(pWrapper config.IWrapper, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(pkg_settings.CServiceName, pR)
		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))

		sett := pWrapper.GetConfig().GetSettings()
		cfgBytes := encoding.Serialize(
			&pkg_config.SConfigSettings{
				FMessageSizeBytes:   sett.GetMessageSizeBytes(),
				FWorkSizeBits:       sett.GetWorkSizeBits(),
				FQueuePeriodMS:      sett.GetQueuePeriodMS(),
				FKeySizeBits:        sett.GetKeySizeBits(),
				FLimitVoidSizeBytes: sett.GetLimitVoidSizeBytes(),
			},
			true,
		)

		api.Response(pW, http.StatusOK, string(cfgBytes))
	}
}
