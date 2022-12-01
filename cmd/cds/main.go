package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/number571/go-peer/modules/crypto/random"
	"github.com/number571/go-peer/modules/storage"
)

func main() {
	if len(os.Args) != 5 {
		panic(fmt.Sprintf(
			"usage: \n\t%s",
			"./main (get|put|del|gen) [path] [storage-password] [data-password]",
		))
	}

	var (
		pathToStorage   = os.Args[2]
		storagePassword = os.Args[3]
		dataPassword    = os.Args[4]
	)

	stg, err := storage.NewCryptoStorage(pathToStorage, []byte(storagePassword), 20)
	if err != nil {
		panic(err)
	}

	switch strings.ToUpper(os.Args[1]) {
	case "GET":
		data, err := stg.Get([]byte(dataPassword))
		if err != nil {
			panic(err)
		}
		fmt.Print(string(data))
	case "PUT":
		if _, err := stg.Get([]byte(dataPassword)); err == nil {
			panic("password already exist")
		}
		err := stg.Set([]byte(dataPassword), readUntilEOF())
		if err != nil {
			panic(err)
		}
	case "DEL":
		err := stg.Del([]byte(dataPassword))
		if err != nil {
			panic(err)
		}
	case "GEN":
		if _, err := stg.Get([]byte(dataPassword)); err == nil {
			panic("password already exist")
		}
		// 1char = 4bit entropy => 128bit
		randStr := random.NewStdPRNG().String(32)
		err := stg.Set([]byte(dataPassword), []byte(randStr))
		if err != nil {
			panic(err)
		}
	default:
		panic("undefined option (get|put)")
	}
}

func readUntilEOF() []byte {
	result := make([]byte, 0, 1024)

	buf := make([]byte, 256)
	file := os.Stdin

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		result = bytes.Join(
			[][]byte{
				result,
				buf[:n],
			},
			[]byte{},
		)
	}

	return result
}
