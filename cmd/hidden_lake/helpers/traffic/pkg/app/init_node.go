package app

import (
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initNode(pCfg config.IConfig, pDBWrapper database.IDBWrapper, pLogger logger.ILogger) network.INode {
	cfgSettings := pCfg.GetSettings()
	queueDuration := time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond

	// queue_period_ms in HLT can be = 0 (as only-storage mode)
	if queueDuration == 0 {
		queueDuration = 1
	}

	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      pCfg.GetAddress().GetTCP(),
			FMaxConnects:  hls_settings.CNetworkMaxConns,
			FReadTimeout:  queueDuration,
			FWriteTimeout: queueDuration,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FNetworkKey:       cfgSettings.GetNetworkKey(),
				FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
				FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
				FLimitVoidSize:    cfgSettings.GetLimitVoidSizeBytes(),
				FWaitReadDeadline: hls_settings.CConnWaitReadDeadline,
				FReadDeadline:     queueDuration,
				FWriteDeadline:    queueDuration,
			}),
		}),
		lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: hls_settings.CNetworkQueueCapacity,
			}),
		),
	).HandleFunc(
		hls_settings.CNetworkMask,
		handler.HandleServiceTCP(pCfg, pDBWrapper, pLogger),
	)
}
