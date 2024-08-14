package main

import (
	"fmt"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/payload"
)

func main() {
	client := client.NewClient(
		client.NewSettings(&client.SSettings{
			FMessageSizeBytes: (8 << 10),
		}),
	)

	key := random.NewCSPRNG().GetBytes(symmetric.CAESKeySize)
	msg, err := client.EncryptMessage(
		key,
		payload.NewPayload64(0x0, []byte("hello, world!")).ToBytes(),
	)
	if err != nil {
		panic(err)
	}

	decMsg, err := client.DecryptMessage(key, msg)
	if err != nil {
		panic(err)
	}

	pld := payload.LoadPayload64(decMsg)
	fmt.Printf("Message: '%s';\n", string(pld.GetBody()))
	fmt.Printf("Encrypted message: '%x'\n", msg)
}
