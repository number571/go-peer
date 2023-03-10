package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/storage"
)

func main() {
	if len(os.Args) != 4 {
		panic(fmt.Sprintf(
			"usage:\n\t%s\nstdin:\n\t%s\n",
			"./main (get|put|del|new) [storage-path] [data-key]",
			"[storage-password]~[data-value]EOF",
		))
	}

	var (
		storagePath = os.Args[2]
		dateKey     = os.Args[3]
	)

	concatData := readUntilEOF()
	splited := bytes.Split(concatData, []byte{'~'})
	if len(splited) == 1 {
		panic("len(splited) == 1")
	}

	storageKey := string(splited[0])
	stg, err := storage.NewCryptoStorage(storagePath, []byte(storageKey), 20)
	if err != nil {
		panic(err)
	}

	switch strings.ToUpper(os.Args[1]) {
	case "GET":
		data, err := stg.Get([]byte(dateKey))
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(data))
	case "PUT":
		if _, err := stg.Get([]byte(dateKey)); err == nil {
			panic("data-key already exist")
		}
		err := stg.Set([]byte(dateKey), bytes.Join(splited[1:], []byte{'~'}))
		if err != nil {
			panic(err)
		}
	case "DEL":
		err := stg.Del([]byte(dateKey))
		if err != nil {
			panic(err)
		}
	case "NEW":
		if _, err := stg.Get([]byte(dateKey)); err == nil {
			panic("data-key already exist")
		}
		// 1char = 4bit entropy => 128bit
		randStr := random.NewStdPRNG().GetString(32)
		err := stg.Set([]byte(dateKey), []byte(randStr))
		if err != nil {
			panic(err)
		}
	default:
		panic("undefined option (get|put|del|new)")
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
