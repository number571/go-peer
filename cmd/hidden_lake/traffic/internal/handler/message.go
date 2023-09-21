package handler

import (
	"io"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
)

func HandleMessageAPI(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger, pNode network.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		httpLogger := http_logger.NewHTTPLogger(hlt_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
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

		switch pR.Method {
		case http.MethodGet:
			query := pR.URL.Query()
			msg, err := database.Load(query.Get("hash"))
			if err != nil {
				pLogger.PushWarn(httpLogger.Get("load_hash"))
				api.Response(pW, http.StatusNotFound, "failed: load message")
				return
			}

			pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, msg.ToString())
			return

		case http.MethodPost:
			msgStringAsBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(httpLogger.Get(http_logger.CLogDecodeBody))
				api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}
			msgString := string(msgStringAsBytes)

			handler := HandleServiceTCP(pCfg, pWrapperDB, pLogger)
			if err := handler(pNode, nil, message.FromStringToBytes(msgString)); err != nil {
				// internal logger
				api.Response(pW, http.StatusBadRequest, "failed: handle message")
				return
			}

			pLogger.PushInfo(httpLogger.Get(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: handle message")
			return
		}
	}
}
