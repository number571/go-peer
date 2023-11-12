package main

import (
	"fmt"
	"strconv"
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
	defer service1.Stop()

	time.Sleep(time.Second) // wait

	if err := service2.AddConnection(serviceAddress); err != nil {
		panic(err)
	}

	msg := message.NewMessage(
		service2.GetSettings().GetConnSettings(),
		payload.NewPayload(
			serviceHeader,
			[]byte("0"),
		),
	)
	service2.BroadcastMessage(msg)

	select {}
}

func handler(serviceName string) network.IHandlerF {
	return func(n network.INode, _ conn.IConn, msg message.IMessage) error {
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
		n.BroadcastMessage(message.NewMessage(
			n.GetSettings().GetConnSettings(),
			payload.NewPayload(
				serviceHeader,
				[]byte(fmt.Sprintf("%d", num+1)),
			),
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
