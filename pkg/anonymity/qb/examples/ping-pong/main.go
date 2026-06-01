package main

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	nodeAddress = "127.0.0.1:8080"
)

func init() {
	printTagVersion()
}

func main() {
	nodeService := runServiceNode()
	nodeClient := runClientNode()

	keyService, _ := exchangeKeys(nodeService, nodeClient)

	numBytes := encoding.Uint64ToBytes(0)
	_ = nodeClient.fAnonymity.SendPayload(
		context.Background(),
		keyService,
		numBytes[:],
	)

	select {}
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
