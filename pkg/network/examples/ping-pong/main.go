package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

var handler = func(id string) network.IHandlerF {
	return func(ctx context.Context, n network.INode, _ conn.IConn, msg layer1.IMessage) error {
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
			layer1.NewMessage(
				layer1.NewConstructSettings(&layer1.SConstructSettings{
					FSettings: n.GetSettings().GetConnSettings().GetMessageSettings(),
				}),
				payload.NewPayload32(serviceHeader, []byte(fmt.Sprintf("%d", num+1))),
			),
		)

		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		_     = runServiceNode(ctx, "node1")
		node1 = runClientNode("node2")
	)

	msg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: node1.GetSettings().GetConnSettings().GetMessageSettings(),
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

func runServiceNode(ctx context.Context, id string) network.INode {
	node := newNode(serviceAddress).HandleFunc(serviceHeader, handler(id))
	go func() { _ = node.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
