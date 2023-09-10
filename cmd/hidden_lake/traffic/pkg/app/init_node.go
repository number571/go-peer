package app

import (
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initNode(pCfg config.IConfig, pWrapperDB database.IWrapperDB, pLogger logger.ILogger) network.INode {
	queueDuration := time.Duration(pCfg.GetSettings().GetQueuePeriodMS()) * time.Millisecond
	if queueDuration == 0 {
		queueDuration = 1 // queue_period_ms in HLT can be = 0 (as only-storage mode)
	}
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      pCfg.GetAddress().GetTCP(),
			FMaxConnects:  hls_settings.CNetworkMaxConns,
			FCapacity:     hls_settings.CNetworkCapacity,
			FWriteTimeout: queueDuration,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FNetworkKey:       pCfg.GetNetworkKey(),
				FMessageSizeBytes: pCfg.GetSettings().GetMessageSizeBytes(),
				FLimitVoidSize:    pCfg.GetSettings().GetLimitVoidSizeBytes(),
				FWaitReadDeadline: hls_settings.CConnWaitReadDeadline,
				FReadDeadline:     queueDuration,
				FWriteDeadline:    queueDuration,
			}),
		}),
	).HandleFunc(
		hls_settings.CNetworkMask,
		handler.HandleServiceTCP(pCfg, pWrapperDB, pLogger),
	)
}
