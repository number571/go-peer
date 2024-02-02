package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf(
			"usage: \n\t%s\n\n",
			"./main [key-size]",
		))
	}

	keySize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if keySize < 0 {
		panic("key size is negative")
	}

	priv := asymmetric.NewRSAPrivKey(uint64(keySize))
	if priv == nil {
		panic("generate key error")
	}

	if err := os.WriteFile("priv.key", []byte(priv.ToString()), 0o644); err != nil {
		panic(err)
	}
	if err := os.WriteFile("pub.key", []byte(priv.GetPubKey().ToString()), 0o644); err != nil {
		panic(err)
	}
}
