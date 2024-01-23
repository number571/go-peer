package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
)

func HandleServiceTCP(pCfg config.IConfig, pDBWrapper database.IDBWrapper, pLogger logger.ILogger) network.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}

	return func(pCtx context.Context, pNode network.INode, pConn conn.IConn, pNetMsg net_message.IMessage) error {
		logBuilder := anon_logger.NewLogBuilder(hlt_settings.CServiceName)

		// enrich logger
		pld := pNetMsg.GetPayload()
		logBuilder.
			WithConn(pConn).
			WithHash(pNetMsg.GetHash()).
			WithProof(pNetMsg.GetProof()).
			WithSize(len(pld.GetBody()))

		if _, err := message.LoadMessage(pCfg.GetSettings(), pld.GetBody()); err != nil {
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			return fmt.Errorf("load message: %w", err)
		}

		// check database
		hltDB := pDBWrapper.Get()
		if hltDB == nil {
			pLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroDatabaseGet))
			return errors.New("database is nil")
		}

		// check message from in database queue
		if err := hltDB.Push(pNetMsg); err != nil {
			if errors.Is(err, database.GErrMessageIsExist) {
				pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoExist))
				return nil
			}
			pLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroDatabaseSet))
			return fmt.Errorf("put message to database: %w", err)
		}

		hasBroadastError := false
		if err := pNode.BroadcastMessage(pCtx, pNetMsg); err != nil {
			// need pass error (some of connections may be closed)
			hasBroadastError = true
		}

		switch hasBroadastError {
		case true:
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
		default:
			pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
		}

		consumers := pCfg.GetConsumers()

		wg := sync.WaitGroup{}
		wg.Add(len(consumers))

		for _, cHost := range consumers {
			go func(cHost string) {
				defer wg.Done()
				_, err := api.Request(
					httpClient,
					http.MethodPost,
					fmt.Sprintf("http://%s", cHost),
					pNetMsg.ToString(),
				)
				if err != nil {
					pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseGetResponse))
					return
				}
				pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogBaseGetResponse))
			}(cHost)
		}

		wg.Wait()

		return nil
	}
}
