package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	nodeAddress = "127.0.0.1:8080"
	nodeRouter  = uint32(0xA557711A)
)

func main() {
	sharedKey := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456")
	nodeService, nodeClient := runServiceNode(), runClientNode()
	exchangeKey(nodeService, nodeClient, sharedKey)

	ctx := context.Background()
	for {
		resp, _ := nodeClient.FetchPayload(
			ctx,
			sharedKey,
			payload.NewPayload32(nodeRouter, []byte("hello, world!")),
		)
		fmt.Println(string(resp))
	}
}

func runClientNode() anonymity.INode {
	ctx := context.Background()
	node := newNode("cnode", "")

	go func() { _ = node.Run(ctx) }()
	node.GetNetworkNode().AddConnection(ctx, nodeAddress)

	return node
}

func runServiceNode() anonymity.INode {
	ctx := context.Background()
	node := newNode("snode", nodeAddress).HandleFunc(
		nodeRouter,
		func(_ context.Context, _ anonymity.INode, _ []byte, b []byte) ([]byte, error) {
			return []byte(fmt.Sprintf("echo: %s", string(b))), nil
		},
	)

	go func() { _ = node.Run(ctx) }()
	go func() { _ = node.GetNetworkNode().Listen(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}

func exchangeKey(node1, node2 anonymity.INode, sharedKey []byte) {
	node1.GetListKeys().AddKey(sharedKey)
	node2.GetListKeys().AddKey(sharedKey)
}
