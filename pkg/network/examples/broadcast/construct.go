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

func newNode(serviceAddress, id string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:     serviceAddress,
			FMaxConnects: 3,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings:       layer1.NewSettings(&layer1.SSettings{}),
				FLimitMessageSizeBytes: (1 << 10),
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
			FWriteTimeout: time.Minute,
			FReadTimeout:  time.Minute,
		}),
		handler(id),
		cache.NewLRUCache(1<<10),
	)
}

var handler = func(serviceName string) network.IHandlerF {
	return func(ctx context.Context, node network.INode, _ conn.IConn, msg layer1.IMessage) error {
		defer func() { _ = node.BroadcastMessage(ctx, msg) }() // send this message to other connections
		fmt.Printf("'%s' got '%s'\n", serviceName, string(msg.GetBody()))
		return nil
	}
}
