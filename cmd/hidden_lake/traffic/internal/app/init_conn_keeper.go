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

func initConnKeeper(cfg config.IConfig, wDB database.IWrapperDB, logger logger.ILogger) conn_keeper.IConnKeeper {
	httpClient := &http.Client{Timeout: time.Minute}
	anonLogger := anon_logger.NewLogger(hlt_settings.CServiceName)
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return []string{cfg.GetConnection()} },
			FDuration:    hls_settings.CNetworkWaitTime,
		}),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FMaxConnects: 1, // one HLS from cfg.Connection()
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.GetNetwork(),
					FMessageSize: hls_settings.CMessageSize,
					FTimeWait:    hls_settings.CNetworkWaitTime,
				}),
			}),
		).HandleFunc(
			hls_settings.CNetworkMask,
			func(_ network.INode, conn conn.IConn, msgBytes []byte) {
				msg := message.LoadMessage(
					msgBytes,
					message.NewParams(
						hls_settings.CMessageSize,
						hls_settings.CWorkSize,
					),
				)
				if msg == nil {
					logger.PushWarn(anonLogger.GetFmtLog(anon_logger.CLogWarnMessageNull, nil, 0, nil, conn))
					return
				}

				var (
					hash  = msg.GetBody().GetHash()
					proof = msg.GetBody().GetProof()
				)

				database := wDB.Get()

				strHash := encoding.HexEncode(hash)
				if _, err := database.Load(strHash); err == nil {
					logger.PushInfo(anonLogger.GetFmtLog(anon_logger.CLogInfoExist, hash, proof, nil, conn))
					return
				}

				if err := database.Push(msg); err != nil {
					logger.PushErro(anonLogger.GetFmtLog(anon_logger.CLogErroDatabaseSet, hash, proof, nil, conn))
					return
				}

				for _, cHost := range cfg.GetConsumers() {
					_, err := api.Request(
						httpClient,
						http.MethodPost,
						fmt.Sprintf("http://%s", cHost),
						msgBytes,
					)
					if err != nil {
						logger.PushWarn(anonLogger.GetFmtLog(anon_logger.CLogWarnUnknownRoute, hash, proof, nil, conn))
						continue
					}
				}

				logger.PushInfo(anonLogger.GetFmtLog(anon_logger.CLogInfoUndecryptable, hash, proof, nil, conn))
			},
		),
	)
}
