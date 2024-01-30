package app

import (
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
)

func initNode(pMutex *sync.Mutex, pCfgW config.IWrapper, pPrivKey asymmetric.IPrivKey, pLogger logger.ILogger, pParallel uint64) anonymity.INode {
	cfg := pCfgW.GetConfig()
	cfgSettings := cfg.GetSettings()
	queueDuration := time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   pkg_settings.CServiceName,
			FF2FDisabled:   cfgSettings.GetF2FDisabled(),
			FNetworkMask:   pkg_settings.CNetworkMask,
			FRetryEnqueue:  pkg_settings.CRetryEnqueue,
			FFetchTimeWait: pkg_settings.CFetchTimeout,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		pLogger,
		anonymity.NewDBWrapper(),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      cfg.GetAddress().GetTCP(),
				FMaxConnects:  pkg_settings.CNetworkMaxConns,
				FReadTimeout:  queueDuration,
				FWriteTimeout: queueDuration,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FNetworkKey:       cfgSettings.GetNetworkKey(),
					FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
					FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
					FLimitVoidSize:    cfgSettings.GetLimitVoidSizeBytes(),
					FWaitReadDeadline: pkg_settings.CConnWaitReadDeadline,
					FReadDeadline:     queueDuration,
					FWriteDeadline:    queueDuration,
				}),
			}),
			lru.NewLRUCache(
				lru.NewSettings(&lru.SSettings{
					FCapacity: pkg_settings.CNetworkQueueCapacity,
				}),
			),
		),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: pkg_settings.CQueueCapacity,
				FVoidCapacity: pkg_settings.CQueuePoolCapacity,
				FParallel:     pParallel,
				FDuration:     queueDuration,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
					FKeySizeBits:      pPrivKey.GetSize(),
				}),
				pPrivKey,
			),
		).WithNetworkSettings(
			pkg_settings.CNetworkMask,
			net_message.NewSettings(&net_message.SSettings{
				FNetworkKey:   cfgSettings.GetNetworkKey(),
				FWorkSizeBits: cfgSettings.GetWorkSizeBits(),
			}),
		),
		func() asymmetric.IListPubKeys {
			f2f := asymmetric.NewListPubKeys()
			for _, pubKey := range cfg.GetFriends() {
				f2f.AddPubKey(pubKey)
			}
			return f2f
		}(),
	).HandleFunc(
		pkg_settings.CServiceMask,
		handler.HandleServiceTCP(pMutex, pCfgW, pLogger),
	)
}
