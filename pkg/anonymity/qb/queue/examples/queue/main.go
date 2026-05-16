package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hybrid/layer1"
	"github.com/number571/go-peer/pkg/crypto/hybrid/macro"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	payloadHead = 0x1
)

func main() {
	privKey := asymmetric.NewPrivKey()
	q := queue.NewQBProblemProcessor(
		queue.NewSettings(&queue.SSettings{
			FMessageConstructSettings: layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: layer1.NewSettings(&layer1.SSettings{}),
			}),
			FQueuePeriod:  time.Second,
			FQueuePoolCap: [2]uint64{1 << 5, 1 << 5},
			FConsumersCap: 1,
		}),
		macro.NewScheme(
			privKey,
			(8<<10),
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
			privKey.GetPubKey(),
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
		if uint64(len(netMsg.GetPayload().GetBody())) != (q.GetScheme().GetMessageSize()) {
			panic("payload size is invalid")
		}
		pubKey, decMsg, err := q.GetScheme().DecryptMessage(
			asymmetric.NewMapPubKeys(privKey.GetPubKey()),
			netMsg.GetPayload().GetBody(),
		)
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
		if !bytes.Equal(pubKey.ToBytes(), privKey.GetPubKey().ToBytes()) {
			panic("public key is invalid")
		}
		fmt.Println(string(pld.GetBody()))
	}
}
