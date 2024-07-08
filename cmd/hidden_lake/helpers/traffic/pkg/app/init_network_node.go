package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initNetworkNode(pDatabase database.IDatabase) {
	cfgSettings := p.fConfig.GetSettings()
	p.fNode = network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      p.fConfig.GetAddress().GetTCP(),
			FMaxConnects:  hls_settings.CNetworkMaxConns,
			FReadTimeout:  hls_settings.CNetworkReadTimeout,
			FWriteTimeout: hls_settings.CNetworkWriteTimeout,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes() + cfgSettings.GetRandMessageSizeBytes(),
				FWorkSizeBits:          cfgSettings.GetWorkSizeBits(),
				FWaitReadTimeout:       hls_settings.CConnWaitReadTimeout,
				FDialTimeout:           hls_settings.CConnDialTimeout,
				FReadTimeout:           hls_settings.CNetworkReadTimeout,
				FWriteTimeout:          hls_settings.CNetworkWriteTimeout,
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
