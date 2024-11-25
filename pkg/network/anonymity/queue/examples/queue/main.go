package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	keySize     = 1024
	payloadHead = 0x1
)

func main() {
	q := queue.NewQBProblemProcessor(
		queue.NewSettings(&queue.SSettings{
			FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{}),
			}),
			FQueuePeriod:  time.Second,
			FPoolCapacity: [2]uint64{1 << 5, 1 << 5},
			FConsumersCap: 1,
		}),
		client.NewClient(
			asymmetric.NewPrivKey(),
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
			q.GetClient().GetPrivKey().GetPubKey(),
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
		msg, err := message.LoadMessage(q.GetClient().GetMessageSize(), netMsg.GetPayload().GetBody())
		if err != nil {
			panic(err)
		}
		pubKey, decMsg, err := q.GetClient().DecryptMessage(
			asymmetric.NewMapPubKeys(q.GetClient().GetPrivKey().GetPubKey()),
			msg.ToBytes(),
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
		if !bytes.Equal(pubKey.ToBytes(), q.GetClient().GetPrivKey().GetDSAPrivKey().GetPubKey().ToBytes()) {
			panic("public key is invalid")
		}
		fmt.Println(string(pld.GetBody()))
	}
}
