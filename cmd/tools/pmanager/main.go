package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func main() {
	saltParam := flag.String("salt", "_salt_", "default salt value")
	workParam := flag.Uint("work", 24, "default work value")
	flag.Parse()

	keyBuilder := keybuilder.NewKeyBuilder(1<<(*workParam), []byte(*saltParam))
	gotPassword := keyBuilder.Build(readUntilEOL(), symmetric.CAESKeySize)

	fmt.Println(base64.URLEncoding.EncodeToString(gotPassword))
}

func readUntilEOL() string {
	res, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		panic(err)
	}
	return string(res)
}
