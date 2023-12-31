package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
)

const (
	workSize = 20
)

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Sprintf(
			"usage:\n\t%s\nstdin:\n\t%s\n\n",
			"./main [service] [login]",
			"[master-key]EOF",
		))
	}

	var (
		service = []byte(os.Args[1])
		login   = []byte(os.Args[2])
	)

	keyBuilder := keybuilder.NewKeyBuilder(1<<workSize, login)
	extendedKey := keyBuilder.Build(readUntilEOF("> "))

	passBytes := hashing.NewHMACSHA256Hasher(extendedKey, service).ToBytes()
	fmt.Println(base64.StdEncoding.EncodeToString(passBytes))
}

func readUntilEOF(prefix string) string {
	fmt.Print(prefix)
	res, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		panic(err)
	}
	return string(res)
}
