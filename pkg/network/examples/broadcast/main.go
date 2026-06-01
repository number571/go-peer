package main

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/network"
)

// {node4} -[msg]-> {node1} -[msg]-> {node2}
//                          -[msg]-> {node3}

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = "127.0.0.1:8080"
)

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
			[]byte("hello, world!"),
		),
	)

	time.Sleep(time.Second)
}

func runClientNode(id string) network.INode {
	ctx := context.Background()
	node := newNode("", id)

	_ = node.AddConnection(ctx, serviceAddress)
	return node
}

func runServiceNode(ctx context.Context, id string) network.INode {
	node := newNode(serviceAddress, id)
	go func() { _ = node.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
