package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_config "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func HandleConfigSettingsAPI(
	pWrapper config.IWrapper,
	pLogger logger.ILogger,
	pNode anonymity.INode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, pkg_config.GetConfigSettings(
			pWrapper.GetConfig(),
			pNode.GetMessageQueue().GetClient(),
		))
	}
}
