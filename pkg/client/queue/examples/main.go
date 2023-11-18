package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	payloadHead = 0x1
)

func main() {
	q := queue.NewMessageQueue(
		queue.NewSettings(&queue.SSettings{
			FDuration:     time.Second,
			FMainCapacity: 1 << 5,
			FPoolCapacity: 1 << 5,
		}),
		client.NewClient(
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: (1 << 12),
			}),
			asymmetric.NewRSAPrivKey(1024),
		),
	)

	if err := q.Run(); err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		msg, err := q.GetClient().EncryptPayload(
			q.GetClient().GetPubKey(),
			payload.NewPayload(payloadHead, []byte(fmt.Sprintf("hello, world! %d", i))),
		)
		if err != nil {
			panic(err)
		}
		if err := q.EnqueueMessage(msg); err != nil {
			panic(err)
		}
	}

	for i := 0; i < 3; i++ {
		msg := <-q.DequeueMessage()
		pubKey, pld, err := q.GetClient().DecryptMessage(msg)
		if err != nil {
			panic(err)
		}
		if pld.GetHead() != payloadHead {
			panic("payload head is invalid")
		}
		if pubKey.GetAddress().ToString() != q.GetClient().GetPubKey().GetAddress().ToString() {
			panic("public key is invalid")
		}
		fmt.Println(string(pld.GetBody()))
	}
}
