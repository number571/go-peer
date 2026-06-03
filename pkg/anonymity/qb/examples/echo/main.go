package main

import (
	"context"
	"fmt"
	"time"
)

const (
	nodeAddress = "127.0.0.1:8080"
	nodeRouter  = uint32(0xA557711A)
)

func init() {
	printTagVersion()
}

func main() {
	nodeService := runServiceNode()
	nodeClient := runClientNode()

	keyToService, _ := exchangeKeys(nodeService, nodeClient)

	ctx := context.Background()
	for {
		resp, _ := nodeClient.fAnonymity.FetchPayload(
			ctx,
			keyToService,
			[]byte("hello, world!"),
		)
		fmt.Println(string(resp))
	}
}

func runClientNode() *sNode {
	ctx := context.Background()
	node := newNode("cnode", "")

	go func() { _ = node.fAnonymity.Run(ctx) }()
	_ = node.fNetwork.AddConnection(ctx, nodeAddress)

	return node
}

func runServiceNode() *sNode {
	ctx := context.Background()
	node := newNode("snode", nodeAddress)

	go func() { _ = node.fAnonymity.Run(ctx) }()
	go func() { _ = node.fNetwork.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
