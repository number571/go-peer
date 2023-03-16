package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/database"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	internal_logger "github.com/number571/go-peer/internal/logger"
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
		internal_logger.StdLogger(cfg.GetLogging()),
		anonymity.NewWrapperDB().Set(database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:    pkg_settings.CPathDB,
				FHashing: true,
			}),
		)),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:     cfg.GetAddress().GetTCP(),
				FCapacity:    pkg_settings.CNetworkCapacity,
				FMaxConnects: pkg_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  cfg.GetNetwork(),
					FMessageSize: pkg_settings.CMessageSize,
					FTimeWait:    pkg_settings.CNetworkWaitTime,
				}),
			}),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     pkg_settings.CQueueCapacity,
				FPullCapacity: pkg_settings.CQueuePullCapacity,
				FDuration:     pkg_settings.CQueueDuration,
			}),
			pkg_settings.InitClient(privKey),
		),
		func() asymmetric.IListPubKeys {
			f2f := asymmetric.NewListPubKeys()
			for _, pubKey := range cfg.GetFriends() {
				f2f.AddPubKey(pubKey)
			}
			return f2f
		}(),
	)
}
