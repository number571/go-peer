package app

import (
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/queue_set"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func initNode(pCfg config.IConfig, pDBWrapper database.IDBWrapper, pLogger logger.ILogger) network.INode {
	cfgSettings := pCfg.GetSettings()

	queueDuration := time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond
	connDeadline := hls_settings.GetConnDeadline(queueDuration)

	// queue_period_ms in HLT can be = 0 (as only-storage mode)
	if queueDuration == 0 {
		queueDuration = 1
		connDeadline = 1
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
				FReadDeadline:     connDeadline,
				FWriteDeadline:    connDeadline,
			}),
		}),
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: hls_settings.CNetworkQueueCapacity,
			}),
		),
	).HandleFunc(
		hls_settings.CNetworkMask,
		handler.HandleServiceTCP(pCfg, pDBWrapper, pLogger),
	)
}
