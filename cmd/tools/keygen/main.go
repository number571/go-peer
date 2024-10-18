package main

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func main() {
	priv := asymmetric.NewPrivKey()
	if err := os.WriteFile("priv.key", []byte(priv.ToString()), 0o600); err != nil {
		panic(err)
	}
	if err := os.WriteFile("pub.key", []byte(priv.GetPubKey().ToString()), 0o600); err != nil {
		panic(err)
	}
}
