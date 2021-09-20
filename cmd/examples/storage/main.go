package main

import (
	"bytes"
	"fmt"

	"github.com/number571/gopeer/crypto"
	"github.com/number571/gopeer/local"
)

const (
	PASSWORD = "privkey-password"
)

func main() {
	priv1 := crypto.NewPrivKey(2048)

	store := local.NewStorage("storage", "storage-password")
	store.Write(priv1, PASSWORD)

	priv2, err := store.Read(PASSWORD)
	if err != nil {
		panic(err)
	}

	fmt.Println(bytes.Equal(priv1.Bytes(), priv2.Bytes()))
}
