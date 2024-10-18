package main

import (
	"flag"
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func main() {
	flag.Parse()
	priv := asymmetric.NewPrivKey()
	if priv == nil {
		panic("generate key error")
	}
	if err := os.WriteFile("priv.key", []byte(priv.ToString()), 0o600); err != nil {
		panic(err)
	}
	if err := os.WriteFile("pub.key", []byte(priv.GetPubKey().ToString()), 0o600); err != nil {
		panic(err)
	}
}
