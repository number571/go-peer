package app

import (
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/queue_set"
)

func (p *sApp) initQueuePusher() {
	p.fQPWrapper.Set(
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: hls_settings.CRequestQueueCapacity,
			}),
		),
	)
}
