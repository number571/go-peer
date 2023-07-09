package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func HandleServiceTCP(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) network.IHandlerF {
	httpClient := &http.Client{Timeout: time.Minute}
	anonLogger := anon_logger.NewLogger(hlt_settings.CServiceName)

	return func(pNode network.INode, pConn conn.IConn, pMsgBytes []byte) {
		msg := message.LoadMessage(
			message.NewSettings(&message.SSettings{
				FMessageSize: hls_settings.CMessageSize,
				FWorkSize:    hls_settings.CWorkSize,
			}),
			pMsgBytes,
		)
		if msg == nil {
			pLogger.PushWarn(anonLogger.GetFmtLog(anon_logger.CLogWarnMessageNull, nil, 0, nil, pConn))
			return
		}

		var (
			hash  = msg.GetBody().GetHash()
			proof = msg.GetBody().GetProof()
		)

		database := pWrapperDB.Get()
		strHash := encoding.HexEncode(hash)
		if _, err := database.Load(strHash); err == nil {
			pLogger.PushInfo(anonLogger.GetFmtLog(anon_logger.CLogInfoExist, hash, proof, nil, pConn))
			return
		}

		if err := database.Push(msg); err != nil {
			pLogger.PushErro(anonLogger.GetFmtLog(anon_logger.CLogErroDatabaseSet, hash, proof, nil, pConn))
			return
		}

		pld := payload.NewPayload(hls_settings.CNetworkMask, pMsgBytes)
		if err := pNode.BroadcastPayload(pld); err != nil {
			pLogger.PushWarn(anonLogger.GetFmtLog(anon_logger.CLogBaseBroadcast, nil, 0, nil, pConn))
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
				pLogger.PushWarn(anonLogger.GetFmtLog(anon_logger.CLogWarnUnknownRoute, hash, proof, nil, pConn))
				continue
			}
		}

		pLogger.PushInfo(anonLogger.GetFmtLog(anon_logger.CLogInfoUndecryptable, hash, proof, nil, pConn))
	}
}
