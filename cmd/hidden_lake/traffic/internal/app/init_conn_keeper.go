package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func initConnKeeper(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) conn_keeper.IConnKeeper {
	httpClient := &http.Client{Timeout: time.Minute}
	anonLogger := anon_logger.NewLogger(hlt_settings.CServiceName)
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return []string{pCfg.GetConnection()} },
			FDuration:    hlt_settings.CTimeWait,
		}),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FMaxConnects:  1, // only one HLS from cfg.Connection()
				FCapacity:     hls_settings.CNetworkCapacity,
				FWriteTimeout: hls_settings.CNetworkWriteTimeout,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:       pCfg.GetNetwork(),
					FMessageSize:      hls_settings.CMessageSize,
					FLimitVoidSize:    hls_settings.CConnLimitVoidSize,
					FWaitReadDeadline: hls_settings.CConnWaitReadDeadline,
					FReadDeadline:     hls_settings.CConnReadDeadline,
					FWriteDeadline:    hls_settings.CConnWriteDeadline,
					FFetchTimeWait:    1, // conn.FetchPayload not used in anonymity package
				}),
			}),
		).HandleFunc(
			hls_settings.CNetworkMask,
			func(_ network.INode, pConn conn.IConn, pMsgBytes []byte) {
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
			},
		),
	)
}
