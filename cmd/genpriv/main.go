// go run main.go -pr -pp -ks 2048
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/number571/gopeer/crypto"
)

var (
	KEY_SIZE uint
	PUSH_PUB bool
	PUSH_RAW bool
)

func init() {
	flag.UintVar(&KEY_SIZE, "ks", 0, "key size")
	flag.BoolVar(&PUSH_PUB, "pp", false, "print public key")
	flag.BoolVar(&PUSH_RAW, "pr", false, "print keys in raw string")
	flag.Parse()
}

func main() {
	priv := crypto.NewPrivKey(KEY_SIZE)
	if priv == nil {
		fmt.Println("error: priv key is nil")
		os.Exit(1)
	}
	if PUSH_RAW {
		fmt.Printf("%X\n", priv.Bytes())
		if PUSH_PUB {
			fmt.Printf("%X\n", priv.PubKey().Bytes())
		}
		return
	}
	fmt.Println(priv.String())
	if PUSH_PUB {
		fmt.Println(priv.PubKey().String())
	}
}
