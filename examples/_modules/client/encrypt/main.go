package main

import (
	"fmt"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

// TODO!!!
func main() {
	var (
		client1 = newClient()
		client2 = newClient()
	)

	msg, err := client1.EncryptPayload(
		client2.GetPubKey(),
		payload.NewPayload(0x0, []byte("hello, world!")),
	)
	if err != nil {
		panic(err)
	}

	pubKey, pld, err := client2.DecryptMessage(msg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message: '%s';\nSender's public key: '%s';\n", string(pld.GetBody()), pubKey.ToString())
	fmt.Printf("Encrypted message: '%s'\n", string(msg.ToBytes()))
}

func newClient() client.IClient {
	return client.NewClient(
		client.NewSettings(&client.SSettings{
			FMessageSize: (1 << 12),
		}),
		asymmetric.NewRSAPrivKey(1024),
	)
}
