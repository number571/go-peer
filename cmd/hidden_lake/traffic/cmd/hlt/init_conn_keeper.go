package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

func initConnKeeper(cfg config.IConfig, db database.IKeyValueDB, logger logger.ILogger) conn_keeper.IConnKeeper {
	anonLogger := anon_logger.NewLogger(hlt_settings.CServiceName)
	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return []string{cfg.Connection()} },
			FDuration:    hlt_settings.CNetworkWaitTime,
		}),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FMaxConnects: hlt_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.Network(),
					FMessageSize: hlt_settings.CMessageSize,
					FTimeWait:    hlt_settings.CNetworkWaitTime,
				}),
			}),
		).Handle(
			hlt_settings.CNetworkMask,
			func(_ network.INode, conn conn.IConn, reqBytes []byte) {
				msg := message.LoadMessage(
					reqBytes,
					message.NewParams(
						db.Settings().GetMessageSize(),
						db.Settings().GetWorkSize(),
					),
				)
				if msg == nil {
					logger.Warn(anonLogger.FmtLog(anon_logger.CLogWarnMessageNull, nil, 0, nil, conn))
					return
				}

				var (
					hash  = msg.Body().Hash()
					proof = msg.Body().Proof()
				)

				strHash := encoding.HexEncode(msg.Body().Hash())
				if _, err := db.Load(strHash); err == nil {
					logger.Info(anonLogger.FmtLog(anon_logger.CLogInfoExist, hash, proof, nil, conn))
					return
				}

				if err := db.Push(msg); err != nil {
					logger.Erro(anonLogger.FmtLog(anon_logger.CLogErroDatabaseSet, hash, proof, nil, conn))
					return
				}

				logger.Info(anonLogger.FmtLog(anon_logger.CLogInfoUnencryptable, hash, proof, nil, conn))
			},
		),
	)
}
