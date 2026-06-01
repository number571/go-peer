package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
)

// client <-> service

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = "127.0.0.1:8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		_       = runServiceNode(ctx)
		conn, _ = conn.Connect(ctx, connSettings(), serviceAddress)
	)

	_ = conn.WriteMessage(
		ctx,
		layer1.NewMessage(
			layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: conn.GetSettings().GetMessageSettings(),
			}),
			[]byte("hello, world!"),
		),
	)

	readCh := make(chan struct{})
	go func() { <-readCh }()

	recvMsg, _ := conn.ReadMessage(ctx, readCh)
	fmt.Println(string(recvMsg.GetBody()))
}

func runServiceNode(ctx context.Context) network.INode {
	node := newNode(serviceAddress)
	go func() { _ = node.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
