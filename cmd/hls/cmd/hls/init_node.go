package main

import (
	"github.com/number571/go-peer/cmd/hls/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"
)

func initNode(cfg config.IConfig, privKey asymmetric.IPrivKey) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FRetryEnqueue: hls_settings.CRetryEnqueue,
			FTimeWait:     hls_settings.CWaitTime,
		}),
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
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    hls_settings.CWorkSize,
					FMessageSize: hls_settings.CMessageSize,
				}),
				privKey,
			),
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
