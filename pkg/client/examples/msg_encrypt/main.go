package main

import (
	"fmt"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

func main() {
	var (
		client1 = newClient()
		client2 = newClient()
	)

	msg, err := client1.EncryptMessage(
		client2.GetPubKey(),
		payload.NewPayload64(0x0, []byte("hello, world!")).ToBytes(),
	)
	if err != nil {
		panic(err)
	}

	pubKey, decMsg, err := client2.DecryptMessage(msg)
	if err != nil {
		panic(err)
	}

	pld := payload.LoadPayload64(decMsg)
	fmt.Printf("Message: '%s';\nSender's public key: '%s';\n", string(pld.GetBody()), pubKey.ToString())
	fmt.Printf("Encrypted message: '%x'\n", msg)
}

func newClient() client.IClient {
	keySize := uint64(1024)
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: (8 << 10),
			FKeySizeBits:      keySize,
		}),
		asymmetric.NewRSAPrivKey(keySize),
	)
}
