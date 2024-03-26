package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

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

	msg := message.NewMessage(
		message.NewSettings(&message.SSettings{
			FNetworkKey:   service2.GetVSettings().GetNetworkKey(),
			FWorkSizeBits: service2.GetSettings().GetConnSettings().GetWorkSizeBits(),
		}),
		payload.NewPayload(
			serviceHeader,
			[]byte("0"),
		),
		1,
		0,
	)
	service2.BroadcastMessage(ctx, msg)

	select {}
}

func handler(serviceName string) network.IHandlerF {
	return func(ctx context.Context, n network.INode, _ conn.IConn, msg message.IMessage) error {
		time.Sleep(time.Second) // delay for view "ping-pong" game

		num, err := strconv.Atoi(string(msg.GetPayload().GetBody()))
		if err != nil {
			return err
		}

		val := "ping"
		if num%2 == 1 {
			val = "pong"
		}

		fmt.Printf("service '%s' got '%s#%d'\n", serviceName, val, num)
		n.BroadcastMessage(
			ctx,
			message.NewMessage(
				message.NewSettings(&message.SSettings{
					FNetworkKey:   n.GetVSettings().GetNetworkKey(),
					FWorkSizeBits: n.GetSettings().GetConnSettings().GetWorkSizeBits(),
				}),
				payload.NewPayload(
					serviceHeader,
					[]byte(fmt.Sprintf("%d", num+1)),
				),
				1,
				0,
			),
		)

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
		lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: (1 << 10),
			}),
		),
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
