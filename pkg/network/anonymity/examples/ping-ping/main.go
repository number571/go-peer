package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	nodeAddress = "127.0.0.1:8080"
	nodeRouter  = uint32(0xA557711A)
)

var (
	handler = func(ctx context.Context, n anonymity.INode, pKey []byte, b []byte) ([]byte, error) {
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
			pKey,
			payload.NewPayload64(uint64(nodeRouter), numBytes[:]),
		)
		return nil, nil
	}
)

func main() {
	sharedKey := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456")
	nodeService, nodeClient := runServiceNode(), runClientNode()
	exchangeKey(nodeService, nodeClient, sharedKey)

	numBytes := encoding.Uint64ToBytes(0)
	_ = nodeClient.SendPayload(
		sharedKey,
		payload.NewPayload64(uint64(nodeRouter), numBytes[:]),
	)

	select {}
}

func runClientNode() anonymity.INode {
	ctx := context.Background()
	node := newNode("cnode", "").HandleFunc(nodeRouter, handler)

	go func() { _ = node.Run(ctx) }()
	node.GetNetworkNode().AddConnection(ctx, nodeAddress)

	return node
}

func runServiceNode() anonymity.INode {
	ctx := context.Background()
	node := newNode("snode", nodeAddress).HandleFunc(nodeRouter, handler)

	go func() { _ = node.Run(ctx) }()
	go func() { _ = node.GetNetworkNode().Listen(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}

func exchangeKey(node1, node2 anonymity.INode, sharedKey []byte) {
	node1.GetListKeys().AddKey(sharedKey)
	node2.GetListKeys().AddKey(sharedKey)
}
