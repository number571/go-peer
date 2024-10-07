package main

import (
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
)

func newNode(serviceAddress string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:     serviceAddress,
			FMaxConnects: 3,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FLimitMessageSizeBytes: (1 << 10),
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
			FWriteTimeout: time.Minute,
			FReadTimeout:  time.Minute,
		}),
		cache.NewLRUCache(1<<10),
	)
}
