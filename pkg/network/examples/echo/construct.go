package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
)

func newNode(serviceAddress string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      serviceAddress,
			FMaxConnects:  2,
			FConnSettings: connSettings(),
			FWriteTimeout: time.Minute,
			FReadTimeout:  time.Minute,
		}),
		handler,
		cache.NewLRUCache(1<<10),
	)
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FMessageSettings:       layer1.NewSettings(&layer1.SSettings{}),
		FLimitMessageSizeBytes: (1 << 10),
		FWaitReadTimeout:       time.Hour,
		FDialTimeout:           time.Minute,
		FReadTimeout:           time.Minute,
		FWriteTimeout:          time.Minute,
	})
}

var handler = func(ctx context.Context, node network.INode, c conn.IConn, msg layer1.IMessage) error {
	resp := fmt.Sprintf("echo: [%s]", string(msg.GetBody()))
	_ = c.WriteMessage(
		ctx,
		layer1.NewMessage(
			layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: node.GetSettings().GetConnSettings().GetMessageSettings(),
			}),
			[]byte(resp),
		),
	)
	return nil
}
