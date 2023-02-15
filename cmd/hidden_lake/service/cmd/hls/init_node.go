package main

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	internal_logger "github.com/number571/go-peer/internal/logger"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/database"
)

func initNode(cfg config.IConfig, privKey asymmetric.IPrivKey) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  pkg_settings.CServiceName,
			FNetworkMask:  pkg_settings.CNetworkMask,
			FRetryEnqueue: pkg_settings.CRetryEnqueue,
			FTimeWait:     pkg_settings.CWaitTime,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		internal_logger.DefaultLogger(cfg.Logging()),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:    pkg_settings.CPathDB,
				FHashing: true,
			}),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    pkg_settings.CNetworkCapacity,
				FMaxConnects: pkg_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.Network(),
					FMessageSize: pkg_settings.CMessageSize,
					FTimeWait:    pkg_settings.CNetworkWaitTime,
				}),
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     pkg_settings.CQueueCapacity,
				FPullCapacity: pkg_settings.CQueuePullCapacity,
				FDuration:     pkg_settings.CQueueDuration,
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
