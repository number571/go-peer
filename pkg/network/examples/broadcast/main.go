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
	"github.com/number571/go-peer/pkg/queue_set"
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

	conn, err := conn.NewConn(connSettings(), serviceAddress)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	err = conn.WriteMessage(ctx, message.NewMessage(
		conn.GetSettings(),
		payload.NewPayload(
			serviceHeader,
			[]byte("hello, world!"),
		),
		1,
	))
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
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: (1 << 10),
			}),
		),
	)
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FWorkSizeBits:     10,
		FMessageSizeBytes: (1 << 10),
		FWaitReadDeadline: time.Hour,
		FReadDeadline:     time.Minute,
		FWriteDeadline:    time.Minute,
	})
}
