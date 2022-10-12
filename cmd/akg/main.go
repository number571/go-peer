package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: ./main [key-size] [inline|pretty]")
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
		fmt.Println(priv.PubKey())
		fmt.Println(priv)
	case "pretty":
		fmt.Println(priv.PubKey().Format())
		fmt.Println(priv.Format())
	}
}
