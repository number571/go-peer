package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

// client <-> service

const (
	serviceHeader  = 0xDEADBEAF
	serviceAddress = "127.0.0.1:8080"
)

var handler = func(ctx context.Context, node network.INode, c conn.IConn, msg message.IMessage) error {
	resp := fmt.Sprintf("echo: [%s]", string(msg.GetPayload().GetBody()))
	_ = c.WriteMessage(
		ctx,
		message.NewMessage(
			message.NewConstructSettings(&message.SConstructSettings{
				FSettings: node.GetSettings().GetConnSettings().GetMessageSettings(),
			}),
			payload.NewPayload32(serviceHeader, []byte(resp)),
		),
	)
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		_       = runServiceNode(ctx)
		conn, _ = conn.Connect(ctx, connSettings(), serviceAddress)
	)

	_ = conn.WriteMessage(
		ctx,
		message.NewMessage(
			message.NewConstructSettings(&message.SConstructSettings{
				FSettings: conn.GetSettings().GetMessageSettings(),
			}),
			payload.NewPayload32(serviceHeader, []byte("hello, world!")),
		),
	)

	readCh := make(chan struct{})
	go func() { <-readCh }()

	recvMsg, _ := conn.ReadMessage(ctx, readCh)
	fmt.Println(string(recvMsg.GetPayload().GetBody()))
}

func runServiceNode(ctx context.Context) network.INode {
	node := newNode(serviceAddress).HandleFunc(serviceHeader, handler)
	go func() { _ = node.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
