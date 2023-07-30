package app

import (
	"time"

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
	queueDuration := time.Duration(pCfg.GetQueuePeriodMS()) * time.Millisecond
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   pkg_settings.CServiceName,
			FNetworkMask:   pkg_settings.CNetworkMask,
			FRetryEnqueue:  pkg_settings.CRetryEnqueue,
			FFetchTimeWait: pkg_settings.CFetchTimeout,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		pLogger,
		anonymity.NewWrapperDB(),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      pCfg.GetAddress().GetTCP(),
				FCapacity:     pkg_settings.CNetworkCapacity,
				FMaxConnects:  pkg_settings.CNetworkMaxConns,
				FWriteTimeout: queueDuration,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:       pCfg.GetNetwork(),
					FMessageSizeBytes: pCfg.GetMessageSizeBytes(),
					FLimitVoidSize:    pCfg.GetLimitVoidSizeBytes(),
					FWaitReadDeadline: pkg_settings.CConnWaitReadDeadline,
					FReadDeadline:     queueDuration,
					FWriteDeadline:    queueDuration,
					FFetchTimeWait:    1, // conn.FetchPayload not used in anonymity package
				}),
			}),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: pkg_settings.CQueueCapacity,
				FPoolCapacity: pkg_settings.CQueuePoolCapacity,
				FDuration:     queueDuration,
			}),
			pkg_settings.InitClient(pCfg, pPrivKey),
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
