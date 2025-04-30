package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/anonymity"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	nodeAddress = "127.0.0.1:8080"
	nodeRouter  = uint32(0xA557711A)
)

var (
	handler = func(ctx context.Context, n anonymity.INode, pubKey asymmetric.IPubKey, b []byte) ([]byte, error) {
		numBytes := [encoding.CSizeUint64]byte{}
		copy(numBytes[:], b)

		num := encoding.BytesToUint64(numBytes)
		msg := "ping"
		if num%2 == 1 {
			msg = "pong"
		}
		fmt.Printf("%s-%d\n", msg, num)

		numBytes = encoding.Uint64ToBytes(num + 1)
		_ = n.SendPayload(
			ctx,
			pubKey,
			payload.NewPayload64(uint64(nodeRouter), numBytes[:]),
		)
		return nil, nil
	}
)

func main() {
	nodeService, nodeClient := runServiceNode(), runClientNode()
	pubKeyService, _ := exchangeKeys(nodeService, nodeClient)

	numBytes := encoding.Uint64ToBytes(0)
	_ = nodeClient.SendPayload(
		context.Background(),
		pubKeyService,
		payload.NewPayload64(uint64(nodeRouter), numBytes[:]),
	)

	select {}
}

func runClientNode() anonymity.INode {
	ctx := context.Background()
	network, node := newNode("cnode", "")
	node.HandleFunc(nodeRouter, handler)

	go func() { _ = node.Run(ctx) }()
	_ = network.AddConnection(ctx, nodeAddress)

	return node
}

func runServiceNode() anonymity.INode {
	ctx := context.Background()
	network, node := newNode("snode", nodeAddress)
	node.HandleFunc(nodeRouter, handler)

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
