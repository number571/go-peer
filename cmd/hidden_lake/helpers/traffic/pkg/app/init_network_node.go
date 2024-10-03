package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache/lru"

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
				FMessageSettings:       cfgSettings,
				FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes() + cfgSettings.GetRandMessageSizeBytes(),
				FWaitReadTimeout:       hls_settings.CConnWaitReadTimeout,
				FDialTimeout:           hls_settings.CConnDialTimeout,
				FReadTimeout:           hls_settings.CNetworkReadTimeout,
				FWriteTimeout:          hls_settings.CNetworkWriteTimeout,
			}),
		}),
		lru.NewLRUCache(hls_settings.CNetworkQueueCapacity),
	).HandleFunc(
		hls_settings.CNetworkMask,
		handler.HandleServiceTCP(p.fConfig, pDatabase, p.fAnonLogger),
	)
}
