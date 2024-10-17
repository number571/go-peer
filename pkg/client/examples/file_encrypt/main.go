package main

import (
	"bytes"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/quantum"
)

func main() {
	client := newClient()
	if err := encrypt(client, "enc_image.jpg", "image.jpg"); err != nil {
		panic(err)
	}
	if err := decrypt(client, "dec_image.jpg", "enc_image.jpg"); err != nil {
		panic(err)
	}
	if !bytes.Equal(fileHash("image.jpg"), fileHash("dec_image.jpg")) {
		panic("decrypt failed")
	}
}

func newClient() client.IClient {
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: (8 << 10),
			FEncKeySizeBytes:  quantum.CCiphertextSize,
		}),
		quantum.NewPrivKeyChain(
			quantum.NewKEMPrivKey(),
			quantum.NewSignerPrivKey(),
		),
	)
}
