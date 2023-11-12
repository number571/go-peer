package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/internal/msgconv"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/network/queue_pusher"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func HandleServiceTCP(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) network.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}

	queuePusher := queue_pusher.NewQueuePusher(
		queue_pusher.NewSettings(&queue_pusher.SSettings{
			FCapacity: hls_settings.CNetworkQueueSize,
		}),
	)

	return func(pNode network.INode, pConn conn.IConn, pMsg net_message.IMessage) error {
		logBuilder := anon_logger.NewLogBuilder(hlt_settings.CServiceName)

		// enrich logger
		logBuilder.
			WithConn(pConn).
			WithSize(len(pMsg.ToBytes()))

		msg := message.LoadMessage(
			pCfg.GetSettings(),
			pMsg.GetPayload().GetBody(),
		)
		if msg == nil {
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			return errors.NewError("message is nil")
		}

		var (
			hash  = msg.GetBody().GetHash()
			proof = msg.GetBody().GetProof()
			hltDB = pWrapperDB.Get()
		)

		// enrich logger
		logBuilder.
			WithHash(hash).
			WithProof(proof)

		// check database
		if hltDB == nil {
			pLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroDatabaseGet))
			return errors.NewError("database is nil")
		}

		// check message from in memory queue
		if ok := queuePusher.Push(hash); !ok {
			pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoExist))
			return nil
		}

		// check message from in database queue
		if err := hltDB.Push(msg); err != nil {
			if errors.HasError(err, &database.SIsExistError{}) {
				pLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoExist))
				return nil
			}
			pLogger.PushErro(logBuilder.WithType(anon_logger.CLogErroDatabaseSet))
			return errors.WrapError(err, "put message to database")
		}

		if err := pNode.BroadcastMessage(pMsg); err != nil {
			pLogger.PushWarn(logBuilder.WithType(anon_logger.CLogBaseBroadcast))
			// need pass error (some of connections may be closed)
		}

		msgString := msgconv.FromBytesToString(msg.ToBytes())
		if msgString == "" {
			panic("got invalid result (func=FromBytesToString)")
		}

		for _, cHost := range pCfg.GetConsumers() {
			_, err := api.Request(
				httpClient,
				http.MethodPost,
				fmt.Sprintf("http://%s", cHost),
				msgString,
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
