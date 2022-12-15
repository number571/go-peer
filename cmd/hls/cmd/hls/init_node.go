package main

import (
	"os"

	"github.com/number571/go-peer/cmd/hls/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
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
			FRetryEnqueue: hls_settings.CRetryEnqueue,
			FTimeWait:     hls_settings.CWaitTime,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		logger.NewLogger(loggerSettings(cfg.Logging())),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FHashing: true,
			}),
			hls_settings.CPathDB,
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    hls_settings.CNetworkCapacity,
				FMaxConnects: hls_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.Network(),
					FMessageSize: hls_settings.CMessageSize,
					FTimeWait:    hls_settings.CNetworkWaitTime,
				}),
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     hls_settings.CQueueCapacity,
				FPullCapacity: hls_settings.CQueuePullCapacity,
				FDuration:     hls_settings.CQueueDuration,
			}),
			hls_settings.InitClient(privKey),
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

func loggerSettings(enabled bool) logger.ISettings {
	if !enabled {
		return logger.NewSettings(&logger.SSettings{})
	}
	return logger.NewSettings(&logger.SSettings{
		FInfo: os.Stdout,
		FWarn: os.Stderr,
		FErro: os.Stderr,
	})
}
