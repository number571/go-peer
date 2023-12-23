package handler

import (
	"net/http"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
)

func HandleHashesAPI(pWrapperDB database.IWrapperDB, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlt_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()
		if database == nil {
			pLogger.PushErro(logBuilder.WithMessage("get_database"))
			api.Response(pW, http.StatusInternalServerError, "failed: get database")
			return
		}

		query := pR.URL.Query()
		id, err := strconv.Atoi(query.Get("id"))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_id"))
			api.Response(pW, http.StatusBadRequest, "failed: get id")
			return
		}

		hash, err := database.Hash(uint64(id))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_hashes"))
			api.Response(pW, http.StatusNotAcceptable, "failed: load size from DB")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, encoding.HexEncode(hash))
	}
}
