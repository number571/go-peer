package main

import (
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func main() {
	client := newClient()
	encryptFile(client, client.GetPubKey(), "image.jpg")
	decryptFile(client, "decrypted_")
}

func newClient() client.IClient {
	privKey := asymmetric.LoadRSAPrivKey(privKeyStr)
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
