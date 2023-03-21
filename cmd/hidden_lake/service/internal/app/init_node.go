package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initNode(pCfg config.IConfig, pPrivKey asymmetric.IPrivKey, pLogger logger.ILogger) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  pkg_settings.CServiceName,
			FNetworkMask:  pkg_settings.CNetworkMask,
			FRetryEnqueue: pkg_settings.CRetryEnqueue,
			FTimeWait:     pkg_settings.CWaitTime,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		pLogger,
		anonymity.NewWrapperDB(),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:     pCfg.GetAddress().GetTCP(),
				FCapacity:    pkg_settings.CNetworkCapacity,
				FMaxConnects: pkg_settings.CNetworkMaxConns,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:  pCfg.GetNetwork(),
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
			pkg_settings.InitClient(pPrivKey),
		),
		func() asymmetric.IListPubKeys {
			f2f := asymmetric.NewListPubKeys()
			for _, pubKey := range pCfg.GetFriends() {
				f2f.AddPubKey(pubKey)
			}
			return f2f
		}(),
	)
}
