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
	q := queue.NewMessageQueueProcessor(
		queue.NewSettings(&queue.SSettings{
			FQueuePeriod:      time.Second,
			FMainPoolCapacity: 1 << 5,
			FRandPoolCapacity: 1 << 5,
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
		err := q.EnqueueMessage(
			q.GetClient().GetPubKey(),
			payload.NewPayload64(payloadHead, []byte(fmt.Sprintf("hello, world! %d", i))).ToBytes(),
		)
		if err != nil {
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
		pubKey, decMsg, err := q.GetClient().DecryptMessage(msg.ToBytes())
		if err != nil {
			panic(err)
		}
		pld := payload.LoadPayload64(decMsg)
		if pld == nil {
			panic("payloas is nil")
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
