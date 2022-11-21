package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/filesystem"
)

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Sprintf(
			"usage: \n\t%s",
			"./main [key-size] [inline|pretty]",
		))
	}

	keySize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if keySize < 0 {
		panic("key size is negative")
	}

	mode := strings.ToLower(os.Args[2])
	if mode != "inline" && mode != "pretty" {
		panic("undefined mode [inline|pretty]")
	}

	priv := asymmetric.NewRSAPrivKey(uint64(keySize))
	if priv == nil {
		panic("generate key error")
	}

	switch mode {
	case "inline":
		filesystem.OpenFile("priv.key").Write([]byte(priv.String()))
		filesystem.OpenFile("pub.key").Write([]byte(priv.PubKey().String()))
	case "pretty":
		filesystem.OpenFile("priv.key").Write([]byte(priv.Format()))
		filesystem.OpenFile("pub.key").Write([]byte(priv.PubKey().Format()))
	}
}
