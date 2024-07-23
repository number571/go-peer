package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache/lru"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

// client <-> service1 <-> service2
func main() {
	var (
		service1 = newNode(serviceAddress)
		service2 = newNode("")
	)
	defer func() {
		if err := service1.Close(); err != nil {
			panic(err)
		}
		if err := service2.Close(); err != nil {
			panic(err)
		}
	}()

	service1.HandleFunc(serviceHeader, handler("#1"))
	service2.HandleFunc(serviceHeader, handler("#2"))

	ctx := context.Background()
	go func() {
		err := service1.Listen(ctx)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()
	time.Sleep(time.Second) // wait

	if err := service2.AddConnection(ctx, serviceAddress); err != nil {
		panic(err)
	}

	conn, err := conn.NewConn(connSettings(), vSettings(), serviceAddress)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	err = conn.WriteMessage(
		ctx,
		message.NewMessage(
			message.NewSettings(&message.SSettings{
				FNetworkKey:   conn.GetVSettings().GetNetworkKey(),
				FWorkSizeBits: conn.GetSettings().GetWorkSizeBits(),
			}),
			payload.NewPayload32(
				serviceHeader,
				[]byte("hello, world!"),
			),
		),
	)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}

func handler(serviceName string) network.IHandlerF {
	return func(ctx context.Context, n network.INode, _ conn.IConn, msg message.IMessage) error {
		defer n.BroadcastMessage(ctx, msg)
		fmt.Printf("service '%s' got '%s'\n", serviceName, string(msg.GetPayload().GetBody()))
		return nil
	}
}

func newNode(serviceAddress string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      serviceAddress,
			FMaxConnects:  2,
			FConnSettings: connSettings(),
			FWriteTimeout: time.Minute,
			FReadTimeout:  time.Minute,
		}),
		vSettings(),
		lru.NewLRUCache(1<<10),
	)
}

func vSettings() conn.IVSettings {
	return conn.NewVSettings(&conn.SVSettings{})
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FLimitMessageSizeBytes: (1 << 10),
		FWaitReadTimeout:       time.Hour,
		FDialTimeout:           time.Minute,
		FReadTimeout:           time.Minute,
		FWriteTimeout:          time.Minute,
	})
}
