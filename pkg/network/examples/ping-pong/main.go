package main

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/network"
)

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = ":8080"
)

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
		[]byte("0"),
	)
	_ = node1.BroadcastMessage(ctx, msg)

	select {}
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
