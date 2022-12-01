package main

import (
	"fmt"

	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/payload"
)

func main() {
	var (
		client1 = newClient()
		client2 = newClient()
	)

	msg, err := client1.Encrypt(
		client2.PubKey(),
		payload.NewPayload(0x0, []byte("hello, world!")),
	)
	if err != nil {
		panic(err)
	}

	pubKey, pld, err := client2.Decrypt(msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message: '%s';\nSender's public key: '%s';\n", string(pld.Body()), pubKey.String())
	fmt.Printf("Encrypted message: '%s'\n", string(msg.Bytes()))
}

func newClient() client.IClient {
	return client.NewClient(
		client.NewSettings(&client.SSettings{
			FMessageSize: (1 << 12),
		}),
		asymmetric.NewRSAPrivKey(1024),
	)
}
