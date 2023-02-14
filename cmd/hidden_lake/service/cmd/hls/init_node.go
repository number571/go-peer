package main

import (
	"os"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/database"
)

func initNode(cfg config.IConfig, privKey asymmetric.IPrivKey) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FNetworkMask:  pkg_settings.CNetworkMask,
			FRetryEnqueue: hls_settings.CRetryEnqueue,
			FTimeWait:     hls_settings.CWaitTime,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		logger.NewLogger(loggerSettings(cfg.Logging())),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:    hls_settings.CPathDB,
				FHashing: true,
			}),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    hls_settings.CNetworkCapacity,
				FMaxConnects: hls_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.Network(),
					FMessageSize: pkg_settings.CMessageSize,
					FTimeWait:    pkg_settings.CNetworkWaitTime,
				}),
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     hls_settings.CQueueCapacity,
				FPullCapacity: hls_settings.CQueuePullCapacity,
				FDuration:     hls_settings.CQueueDuration,
			}),
			pkg_settings.InitClient(privKey),
		),
		func() friends.IF2F {
			f2f := friends.NewF2F()
			for _, pubKey := range cfg.Friends() {
				f2f.Append(pubKey)
			}
			return f2f
		}(),
	)
}

func loggerSettings(logging config.ILogging) logger.ISettings {
	sett := &logger.SSettings{}
	if logging.Info() {
		sett.FInfo = os.Stdout
	}
	if logging.Warn() {
		sett.FInfo = os.Stderr
	}
	if logging.Erro() {
		sett.FInfo = os.Stderr
	}
	return logger.NewSettings(sett)
}
