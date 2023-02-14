package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	payloadHead = 0x1
)

func main() {
	q := queue.NewQueue(
		queue.NewSettings(&queue.SSettings{
			FDuration:     time.Second,
			FCapacity:     1 << 5,
			FPullCapacity: 1 << 5,
		}),
		client.NewClient(
			client.NewSettings(&client.SSettings{
				FMessageSize: 1 << 12,
			}),
			asymmetric.NewRSAPrivKey(1024),
		),
	)

	if err := q.Run(); err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		msg, err := q.Client().Encrypt(
			q.Client().PubKey(),
			payload.NewPayload(payloadHead, []byte(fmt.Sprintf("hello, world! %d", i))),
		)
		if err != nil {
			panic(err)
		}
		if err := q.Enqueue(msg); err != nil {
			panic(err)
		}
	}

	for i := 0; i < 3; i++ {
		msg := <-q.Dequeue()
		pubKey, pld, err := q.Client().Decrypt(msg)
		if err != nil {
			panic(err)
		}
		if pld.Head() != payloadHead {
			panic("payload head is invalid")
		}
		if pubKey.Address().String() != q.Client().PubKey().Address().String() {
			panic("public key is invalid")
		}
		fmt.Println(string(pld.Body()))
	}
}
