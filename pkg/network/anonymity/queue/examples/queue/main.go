package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	keySize     = 1024
	payloadHead = 0x1
)

func main() {
	q := queue.NewMessageQueue(
		queue.NewSettings(&queue.SSettings{
			FDuration:     time.Second,
			FParallel:     1,
			FMainCapacity: 1 << 5,
			FVoidCapacity: 1 << 5,
		}),
		queue.NewVSettings(&queue.SVSettings{}),
		client.NewClient(
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: (1 << 12),
				FKeySizeBits:      keySize,
			}),
			asymmetric.NewRSAPrivKey(1024),
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := q.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	for i := 0; i < 3; i++ {
		msg, err := q.GetClient().EncryptPayload(
			q.GetClient().GetPubKey(),
			payload.NewPayload64(payloadHead, []byte(fmt.Sprintf("hello, world! %d", i))),
		)
		if err != nil {
			panic(err)
		}
		if err := q.EnqueueMessage(msg); err != nil {
			panic(err)
		}
	}

	for i := 0; i < 3; i++ {
		netMsg := q.DequeueMessage(ctx)
		if netMsg == nil {
			panic("net message is nil")
		}
		msg, err := message.LoadMessage(q.GetClient().GetSettings(), netMsg.GetPayload().GetBody())
		if err != nil {
			panic(err)
		}
		pubKey, pld, err := q.GetClient().DecryptMessage(msg)
		if err != nil {
			panic(err)
		}
		if pld.GetHead() != payloadHead {
			panic("payload head is invalid")
		}
		if pubKey.GetHasher().ToString() != q.GetClient().GetPubKey().GetHasher().ToString() {
			panic("public key is invalid")
		}
		fmt.Println(string(pld.GetBody()))
	}
}
