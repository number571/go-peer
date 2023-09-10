package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/network/anonymity/logbuilder"
)

func HandleServiceTCP(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) network.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}

	return func(pNode network.INode, pConn conn.IConn, pMsgBytes []byte) {
		logBuilder := logbuilder.NewLogBuilder(hlt_settings.CServiceName)

		// enrich logger
		logBuilder.WithConn(pConn)
		logBuilder.WithSize(len(pMsgBytes))

		msg := message.LoadMessage(
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: pCfg.GetSettings().GetMessageSizeBytes(),
				FWorkSizeBits:     pCfg.GetSettings().GetWorkSizeBits(),
			}),
			pMsgBytes,
		)
		if msg == nil {
			pLogger.PushWarn(logBuilder.Get(logbuilder.CLogWarnMessageNull))
			return
		}

		var (
			hash     = msg.GetBody().GetHash()
			proof    = msg.GetBody().GetProof()
			database = pWrapperDB.Get()
		)

		// enrich logger
		logBuilder.WithHash(hash).WithProof(proof)

		// check/push hash of message
		if db := database.GetOriginal(); db != nil {
			hashDB := []byte(fmt.Sprintf("_hash_%X", hash))
			if _, err := db.Get(hashDB); err == nil {
				pLogger.PushInfo(logBuilder.Get(logbuilder.CLogInfoExist))
				return
			}
			if err := db.Set(hashDB, []byte{1}); err != nil {
				pLogger.PushErro(logBuilder.Get(logbuilder.CLogErroDatabaseSet))
				return
			}
		}

		// push message
		if err := database.Push(msg); err != nil {
			pLogger.PushErro(logBuilder.Get(logbuilder.CLogErroDatabaseSet))
			return
		}

		pld := payload.NewPayload(hls_settings.CNetworkMask, pMsgBytes)
		if err := pNode.BroadcastPayload(pld); err != nil {
			pLogger.PushWarn(logBuilder.Get(logbuilder.CLogBaseBroadcast))
			// need pass error (some of connections may be closed)
		}

		for _, cHost := range pCfg.GetConsumers() {
			_, err := api.Request(
				httpClient,
				http.MethodPost,
				fmt.Sprintf("http://%s", cHost),
				pMsgBytes,
			)
			if err != nil {
				pLogger.PushWarn(logBuilder.Get(logbuilder.CLogWarnUnknownRoute))
				continue
			}
		}

		pLogger.PushInfo(logBuilder.Get(logbuilder.CLogBaseBroadcast))
	}
}
