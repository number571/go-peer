package main

import (
	"bytes"

	"github.com/number571/go-peer/pkg/client"
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
	return client.NewClient(
		asymmetric.NewPrivKey(),
		(10 << 10),
	)
}
