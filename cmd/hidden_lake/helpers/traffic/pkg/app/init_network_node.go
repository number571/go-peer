package app

import (
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initNetworkNode(pDatabase database.IDatabase) {
	cfgSettings := p.fConfig.GetSettings()
	maxQueuePeriod := time.Duration(cfgSettings.GetMaxQueuePeriodMS()) * time.Millisecond

	// queue_period_ms in HLT can be = 0 (as only-storage mode)
	if maxQueuePeriod == 0 {
		maxQueuePeriod = 1
	}

	p.fNode = network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      p.fConfig.GetAddress().GetTCP(),
			FMaxConnects:  hls_settings.CNetworkMaxConns,
			FReadTimeout:  maxQueuePeriod,
			FWriteTimeout: maxQueuePeriod,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FWorkSizeBits:          cfgSettings.GetWorkSizeBits(),
				FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
				FLimitVoidSizeBytes:    cfgSettings.GetLimitVoidSizeBytes(),
				FWaitReadTimeout:       hls_settings.CConnWaitReadTimeout,
				FDialTimeout:           hls_settings.CConnDialTimeout,
				FReadTimeout:           maxQueuePeriod,
				FWriteTimeout:          maxQueuePeriod,
			}),
		}),
		conn.NewVSettings(&conn.SVSettings{
			FNetworkKey: cfgSettings.GetNetworkKey(),
		}),
		lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: hls_settings.CNetworkQueueCapacity,
			}),
		),
	).HandleFunc(
		hls_settings.CNetworkMask,
		handler.HandleServiceTCP(p.fConfig, pDatabase, p.fAnonLogger),
	)
}
