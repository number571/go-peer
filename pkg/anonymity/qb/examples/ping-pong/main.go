package main

import (
	"context"
	"fmt"
	"time"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/crypto/hybrid"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	nodeAddress = "127.0.0.1:8080"
	nodeRouter  = uint32(0xA557711A)
)

var (
	handler = func(ctx context.Context, n anonymity.INode, pKey hybrid.IParticipantKey, b []byte) ([]byte, error) {
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
			pKey,
			payload.NewPayload64(uint64(nodeRouter), numBytes[:]),
		)
		return nil, nil
	}
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
		payload.NewPayload64(uint64(nodeRouter), numBytes[:]),
	)

	select {}
}

func runClientNode() *sNode {
	ctx := context.Background()
	node := newNode("cnode", "")
	node.fAnonymity.HandleFunc(nodeRouter, handler)

	go func() { _ = node.fAnonymity.Run(ctx) }()
	_ = node.fNetwork.AddConnection(ctx, nodeAddress)

	return node
}

func runServiceNode() *sNode {
	ctx := context.Background()
	node := newNode("snode", nodeAddress)
	node.fAnonymity.HandleFunc(nodeRouter, handler)

	go func() { _ = node.fAnonymity.Run(ctx) }()
	go func() { _ = node.fNetwork.Run(ctx) }()

	time.Sleep(time.Second) // wait listener
	return node
}
