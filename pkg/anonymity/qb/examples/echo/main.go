package main

import (
	"context"
	"fmt"
	"time"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	nodeAddress = "127.0.0.1:8080"
	nodeRouter  = uint32(0xA557711A)
)

func main() {
	nodeService, nodeClient := runServiceNode(), runClientNode()
	pubKeyService, _ := exchangeKeys(nodeService, nodeClient)

	ctx := context.Background()
	for {
		resp, _ := nodeClient.FetchPayload(
			ctx,
			pubKeyService,
			payload.NewPayload32(nodeRouter, []byte("hello, world!")),
		)
		fmt.Println(string(resp))
	}
}

func runClientNode() anonymity.INode {
	ctx := context.Background()
	network, node := newNode("cnode", "")

	go func() { _ = node.Run(ctx) }()
	_ = network.AddConnection(ctx, nodeAddress)

	return node
}

func runServiceNode() anonymity.INode {
	ctx := context.Background()
	network, node := newNode("snode", nodeAddress)
	node.HandleFunc(
		nodeRouter,
		func(_ context.Context, _ anonymity.INode, _ asymmetric.IPubKey, b []byte) ([]byte, error) {
			return []byte("echo: " + string(b)), nil
		},
	)

	go func() { _ = node.Run(ctx) }()
	go func() { _ = network.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}

func exchangeKeys(node1, node2 anonymity.INode) (asymmetric.IPubKey, asymmetric.IPubKey) {
	pubKey1 := node1.GetQBProcessor().GetClient().GetPrivKey().GetPubKey()
	pubKey2 := node2.GetQBProcessor().GetClient().GetPrivKey().GetPubKey()

	node1.GetMapPubKeys().SetPubKey(pubKey2)
	node2.GetMapPubKeys().SetPubKey(pubKey1)

	return pubKey1, pubKey2
}
