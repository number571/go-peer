package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func HandleServiceTCP(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) network.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}

	return func(pNode network.INode, pConn conn.IConn, pNetMsg net_message.IMessage) error {
		logBuilder := anon_logger.NewLogBuilder(hlt_settings.CServiceName)

		// enrich logger
		logBuilder.
			WithConn(pConn).
			WithSize(len(pNetMsg.ToBytes()))

		_, err := message.LoadMessage(
			pCfg.GetSettings(),
			pNetMsg.GetPayload().GetBody(),
		)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			return fmt.Errorf("load message: %w", err)
		}

		var (
			hash  = pNetMsg.GetHash()
			proof = pNetMsg.GetProof()
			hltDB = pWrapperDB.Get()
		)

		// enrich logger
		logBuilder.
			WithHash(hash).
			WithProof(proof)

		// check database
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

		if err := pNode.BroadcastMessage(pNetMsg); err != nil {
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
			// need pass error (some of connections may be closed)
		}

		// msgString := msgconv.FromBytesToString(msg.ToBytes())
		// if msgString == "" {
		// 	panic("got invalid result (func=FromBytesToString)")
		// }

		for _, cHost := range pCfg.GetConsumers() {
			_, err := api.Request(
				httpClient,
				http.MethodPost,
				fmt.Sprintf("http://%s", cHost),
				pNetMsg.ToString(),
			)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnUnknownRoute))
				continue
			}
		}

		pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
		return nil
	}
}
