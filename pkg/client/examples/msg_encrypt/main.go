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

	msg, err := client1.EncryptMessage(
		client2.GetPrivKey().GetKEncPrivKey().GetPubKey(),
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
	fmt.Printf("Message: '%s';\nSender's public key: '%X';\n", string(pld.GetBody()), pubKey.ToBytes())
	fmt.Printf("Encrypted message: '%x'\n", msg)
}

func newClient() client.IClient {
	return client.NewClient(
		asymmetric.NewPrivKey(),
		(8 << 10),
	)
}
