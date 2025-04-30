package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
)

// {node4} -[msg]-> {node1} -[msg]-> {node2}
//                          -[msg]-> {node3}

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = "127.0.0.1:8080"
)

var handler = func(serviceName string) network.IHandlerF {
	return func(ctx context.Context, node network.INode, _ conn.IConn, msg layer1.IMessage) error {
		defer func() { _ = node.BroadcastMessage(ctx, msg) }() // send this message to other connections
		fmt.Printf("'%s' got '%s'\n", serviceName, string(msg.GetPayload().GetBody()))
		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		_     = runServiceNode(ctx, "node1")
		_     = runClientNode("node2")
		_     = runClientNode("node3")
		node4 = runClientNode("node4")
	)

	_ = node4.BroadcastMessage(
		context.Background(),
		layer1.NewMessage(
			layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: node4.GetSettings().GetConnSettings().GetMessageSettings(),
			}),
			payload.NewPayload32(
				serviceHeader,
				[]byte("hello, world!"),
			),
		),
	)

	time.Sleep(time.Second)
}

func runClientNode(id string) network.INode {
	ctx := context.Background()
	node := newNode("").HandleFunc(serviceHeader, handler(id))

	_ = node.AddConnection(ctx, serviceAddress)
	return node
}

func runServiceNode(ctx context.Context, id string) network.INode {
	node := newNode(serviceAddress).HandleFunc(serviceHeader, handler(id))
	go func() { _ = node.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
