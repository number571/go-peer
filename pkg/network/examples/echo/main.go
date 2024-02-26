package main

import (
	"context"
	"errors"
	"fmt"
	"net"
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

// client <-> service
func main() {
	var (
		service = newNode(serviceAddress)
	)
	defer func() {
		if err := service.Close(); err != nil {
			panic(err)
		}
	}()

	service.HandleFunc(serviceHeader, handler())

	ctx := context.Background()
	go func() {
		err := service.Listen(ctx)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()
	time.Sleep(time.Second) // wait

	conn, err := conn.NewConn(connSettings(), serviceAddress)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()

	sendMsg := message.NewMessage(
		conn.GetSettings(),
		payload.NewPayload(serviceHeader, []byte("hello, world!")),
		1,
	)
	if err := conn.WriteMessage(ctx, sendMsg); err != nil {
		panic(err)
	}

	readCh := make(chan struct{})
	go func() { <-readCh }()

	recvMsg, err := conn.ReadMessage(ctx, readCh)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(recvMsg.GetPayload().GetBody()))
}

func handler() network.IHandlerF {
	return func(ctx context.Context, node network.INode, c conn.IConn, msg message.IMessage) error {
		c.WriteMessage(ctx, message.NewMessage(
			node.GetSettings().GetConnSettings(),
			payload.NewPayload(
				serviceHeader,
				[]byte(fmt.Sprintf("echo: [%s]", string(msg.GetPayload().GetBody()))),
			),
			1,
		))
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
		lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: (1 << 10),
			}),
		),
	)
}

func connSettings() conn.ISettings {
	return conn.NewSettings(&conn.SSettings{
		FWorkSizeBits:     10,
		FMessageSizeBytes: (1 << 10),
		FWaitReadTimeout:  time.Hour,
		FDialTimeout:      time.Minute,
		FReadTimeout:      time.Minute,
		FWriteTimeout:     time.Minute,
	})
}
