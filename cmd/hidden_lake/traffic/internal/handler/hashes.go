package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

func HandleHashesAPI(pWrapperDB database.IWrapperDB, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlt_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(httpLogger.Get(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()
		if database == nil {
			pLogger.PushErro(httpLogger.Get("get_database"))
			api.Response(pW, http.StatusInternalServerError, "failed: get database")
			return
		}

		hashes, err := database.Hashes()
		if err != nil {
			pLogger.PushErro(httpLogger.Get("get_hashes"))
			api.Response(pW, http.StatusInternalServerError, "failed: load size from DB")
			return
		}

		pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, hashes)
	}
}
