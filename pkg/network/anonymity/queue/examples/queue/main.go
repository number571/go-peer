package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
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
			client.NewSettings(&client.SSettings{
				FMessageSizeBytes: (1 << 12),
			}),
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := q.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	key := random.NewCSPRNG().GetBytes(symmetric.CAESKeySize)
	for i := 0; i < 3; i++ {
		err := q.EnqueueMessage(
			key,
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
		msg := netMsg.GetPayload().GetBody()
		if !q.GetClient().MessageIsValid(msg) {
			panic("message is invalid")
		}
		decMsg, err := q.GetClient().DecryptMessage(key, msg)
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
		fmt.Println(string(pld.GetBody()))
	}
}
