package main

import (
	"bytes"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
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
	privKey := asymmetric.NewRSAPrivKey(1024)
	return client.NewClient(
		newSettings(privKey.GetSize()),
		privKey,
	)
}

func newSettings(size uint64) message.ISettings {
	return message.NewSettings(&message.SSettings{
		FMessageSizeBytes: (8 << 10),
		FKeySizeBits:      size,
	})
}
