package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/stringtools"
)

func HandleConfigConnectsAPI(pWrapper config.IWrapper, pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		isBackupParam := pR.URL.Query().Get("is_backup")
		switch isBackupParam {
		case "false", "true":
			// pass
		default:
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_param"))
			api.Response(pW, http.StatusNotAcceptable, "failed: incorrect param")
			return
		}

		isBackup := isBackupParam == "true"

		switch pR.Method {
		case http.MethodGet:
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			if isBackup {
				api.Response(pW, http.StatusOK, pWrapper.GetConfig().GetBackupConnections())
				return
			}
			api.Response(pW, http.StatusOK, pWrapper.GetConfig().GetConnections())
			return
		}

		connectBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			api.Response(pW, http.StatusConflict, "failed: read connect bytes")
			return
		}

		connect := strings.TrimSpace(string(connectBytes))
		if connect == "" {
			pLogger.PushWarn(logBuilder.WithMessage("read_connect"))
			api.Response(pW, http.StatusTeapot, "failed: connect is nil")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			if isBackup {
				connects := stringtools.UniqAppendToSlice(
					pWrapper.GetConfig().GetBackupConnections(),
					connect,
				)
				if err := pWrapper.GetEditor().UpdateBackupConnections(connects); err != nil {
					pLogger.PushWarn(logBuilder.WithMessage("update_backup_connections"))
					api.Response(pW, http.StatusInternalServerError, "failed: update backup connections")
					return
				}

				pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
				api.Response(pW, http.StatusOK, "success: update backup connections")
				return
			}

			connects := stringtools.UniqAppendToSlice(
				pWrapper.GetConfig().GetConnections(),
				connect,
			)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				api.Response(pW, http.StatusInternalServerError, "failed: update connections")
				return
			}

			_ = pNode.GetNetworkNode().AddConnection(connect) // connection may be refused (closed)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: update connections")
			return

		case http.MethodDelete:
			if isBackup {
				connects := stringtools.DeleteFromSlice(pWrapper.GetConfig().GetBackupConnections(), connect)
				if err := pWrapper.GetEditor().UpdateBackupConnections(connects); err != nil {
					pLogger.PushWarn(logBuilder.WithMessage("delete_backup_connections"))
					api.Response(pW, http.StatusInternalServerError, "failed: delete backup connection")
					return
				}

				pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
				api.Response(pW, http.StatusOK, "success: delete backup connection")
				return
			}

			connects := stringtools.DeleteFromSlice(pWrapper.GetConfig().GetConnections(), connect)
			if err := pWrapper.GetEditor().UpdateConnections(connects); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_connections"))
				api.Response(pW, http.StatusInternalServerError, "failed: delete connection")
				return
			}

			if err := pNode.GetNetworkNode().DelConnection(connect); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("del_connections"))
				api.Response(pW, http.StatusInternalServerError, "failed: del connection")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: delete connection")
		}
	}
}
