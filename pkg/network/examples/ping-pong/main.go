package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

var handler = func(id string) network.IHandlerF {
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

		fmt.Printf("'%s' got '%s#%d'\n", id, val, num)
		n.BroadcastMessage(
			ctx,
			message.NewMessage(
				message.NewConstructSettings(&message.SConstructSettings{
					FSettings: n.GetSettings().GetConnSettings(),
				}),
				payload.NewPayload32(serviceHeader, []byte(fmt.Sprintf("%d", num+1))),
			),
		)

		return nil
	}
}

func main() {
	var (
		_     = runServiceNode("node1")
		node1 = runClientNode("node2")
	)

	ctx := context.Background()
	msg := message.NewMessage(
		message.NewConstructSettings(&message.SConstructSettings{
			FSettings: node1.GetSettings().GetConnSettings(),
		}),
		payload.NewPayload32(serviceHeader, []byte("0")),
	)
	node1.BroadcastMessage(ctx, msg)

	select {}
}

func runClientNode(id string) network.INode {
	ctx := context.Background()
	node := newNode("").HandleFunc(serviceHeader, handler(id))

	node.AddConnection(ctx, serviceAddress)
	return node
}

func runServiceNode(id string) network.INode {
	ctx := context.Background()
	node := newNode(serviceAddress).HandleFunc(serviceHeader, handler(id))

	go func() { _ = node.Listen(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
