package main

import (
	"bytes"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hybrid/client"
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
	fmt.Println("done")
}

func newClient() client.IClient {
	return client.NewClient(
		asymmetric.NewPrivKey(),
		(8 << 10),
	)
}
