package main

import (
	"context"
	"fmt"
	"time"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/payload"
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
			payload.NewPayload32(nodeRouter, []byte("hello, world!")),
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
	node.fAnonymity.HandleFunc(
		nodeRouter,
		func(_ context.Context, _ anonymity.INode, _ layer2.IParticipantKey, b []byte) ([]byte, error) {
			return []byte("echo: " + string(b)), nil
		},
	)

	go func() { _ = node.fAnonymity.Run(ctx) }()
	go func() { _ = node.fNetwork.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
