package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/filesystem"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf(
			"usage: \n\t%s",
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

	filesystem.OpenFile("priv.key").Write([]byte(priv.ToString()))
	filesystem.OpenFile("pub.key").Write([]byte(priv.PubKey().ToString()))
}
