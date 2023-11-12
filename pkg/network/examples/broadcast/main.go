package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/network/queue_pusher"
	"github.com/number571/go-peer/pkg/payload"
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

	service1.HandleFunc(serviceHeader, handler("#1"))
	service2.HandleFunc(serviceHeader, handler("#2"))

	if err := service1.Run(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second) // wait

	if err := service2.AddConnection(serviceAddress); err != nil {
		panic(err)
	}

	conn, err := conn.NewConn(connSettings(), serviceAddress)
	if err != nil {
		panic(err)
	}

	err = conn.WriteMessage(message.NewMessage(
		conn.GetSettings(),
		payload.NewPayload(
			serviceHeader,
			[]byte("hello, world!"),
		),
	))
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}

func handler(serviceName string) network.IHandlerF {
	return func(n network.INode, _ conn.IConn, msg message.IMessage) error {
		defer n.BroadcastMessage(msg)
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
		queue_pusher.NewQueuePusher(
			queue_pusher.NewSettings(&queue_pusher.SSettings{
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
