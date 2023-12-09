package main

import (
	"fmt"
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

	if err := service.Listen(); err != nil {
		panic(err)
	}
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
	)
	if err := conn.WriteMessage(sendMsg); err != nil {
		panic(err)
	}

	readCh := make(chan struct{})
	go func() { <-readCh }()

	recvMsg, err := conn.ReadMessage(readCh)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(recvMsg.GetPayload().GetBody()))
}

func handler() network.IHandlerF {
	return func(node network.INode, c conn.IConn, msg message.IMessage) error {
		c.WriteMessage(message.NewMessage(
			node.GetSettings().GetConnSettings(),
			payload.NewPayload(
				serviceHeader,
				[]byte(fmt.Sprintf("echo: [%s]", string(msg.GetPayload().GetBody()))),
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
