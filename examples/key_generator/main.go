package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

func main() {
	if len(os.Args) != 2 {
		panic("usage: ./main [key-size]")
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

	fmt.Println(priv.PubKey().String())
	fmt.Println(priv.String())
}
