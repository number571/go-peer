package app

import (
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initNode(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      pCfg.GetAddress().GetTCP(),
			FMaxConnects:  hls_settings.CNetworkMaxConns,
			FCapacity:     hls_settings.CNetworkCapacity,
			FWriteTimeout: hls_settings.CNetworkWriteTimeout,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FNetworkKey:       pCfg.GetNetwork(),
				FMessageSize:      pCfg.GetMessageSize(),
				FLimitVoidSize:    hls_settings.CConnLimitVoidSize,
				FWaitReadDeadline: hls_settings.CConnWaitReadDeadline,
				FReadDeadline:     hls_settings.CConnReadDeadline,
				FWriteDeadline:    hls_settings.CConnWriteDeadline,
				FFetchTimeWait:    1, // conn.FetchPayload not used in anonymity package
			}),
		}),
	).HandleFunc(
		hls_settings.CNetworkMask,
		handler.HandleServiceTCP(pCfg, pWrapperDB, pLogger),
	)
}
