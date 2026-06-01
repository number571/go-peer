package main

import (
	"context"
	"fmt"
	"strconv"
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
			FMaxConnects: 2,
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

var handler = func(id string) network.IHandlerF {
	return func(ctx context.Context, n network.INode, _ conn.IConn, msg layer1.IMessage) error {
		time.Sleep(time.Second) // delay for view "ping-pong" game

		num, err := strconv.Atoi(string(msg.GetBody()))
		if err != nil {
			return err
		}

		val := "ping"
		if num%2 == 1 {
			val = "pong"
		}

		fmt.Printf("'%s' got '%s#%d'\n", id, val, num)
		_ = n.BroadcastMessage(
			ctx,
			layer1.NewMessage(
				layer1.NewConstructSettings(&layer1.SConstructSettings{
					FSettings: n.GetSettings().GetConnSettings().GetMessageSettings(),
				}),
				[]byte(strconv.Itoa(num+1)),
			),
		)

		return nil
	}
}
