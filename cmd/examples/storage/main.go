package main

import (
	"bytes"
	"fmt"

	"github.com/number571/gopeer/crypto"
	"github.com/number571/gopeer/local"
)

const (
	SUBJECT  = "application_#1"
	PASSWORD = "privkey-password"
)

func main() {
	secret1 := []byte(crypto.NewPrivKey(512).Bytes())

	store := local.NewStorage("storage.enc", "storage-password")
	store.Write(secret1, SUBJECT, PASSWORD)

	secret2, err := store.Read(SUBJECT, PASSWORD)
	if err != nil {
		panic(err)
	}

	fmt.Println(bytes.Equal(secret1, secret2))
}
