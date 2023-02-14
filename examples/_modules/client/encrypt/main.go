package main

import (
	"fmt"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
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

	// fmt.Printf("Encrypted message: '%s'\nHash: '%s'\n", encoding.HexEncode(msg.Bytes()), encoding.HexEncode(msg.Body().Hash()))
}

func newClient() client.IClient {
	return client.NewClient(
		client.NewSettings(&client.SSettings{
			FMessageSize: (1 << 12),
		}),
		asymmetric.NewRSAPrivKey(4096),
	)
}

// func newClient() client.IClient {
// 	return hls_settings.InitClient(asymmetric.NewRSAPrivKey(hls_settings.CAKeySize))
// }
