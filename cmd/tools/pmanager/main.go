package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/number571/go-peer/pkg/crypto/entropy"
	"github.com/number571/go-peer/pkg/crypto/hashing"
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

	booster := entropy.NewEntropyBooster(workSize, login)
	extendedKey := booster.BoostEntropy(readUntilEOF())

	passBytes := hashing.NewHMACSHA256Hasher(extendedKey, service).ToBytes()
	fmt.Println(base64.StdEncoding.EncodeToString(passBytes))
}

func readUntilEOF() []byte {
	res, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return res
}
