package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/number571/go-peer/pkg/storage"
)

const (
	cWorkSize = 20 // bits
)

func main() {
	if len(os.Args) != 4 {
		panic(fmt.Sprintf(
			"usage:\n\t%s\nstdin:\n\t%s\n\t%s\n\n",
			"./main (get|put|del) [storage-path] [data-key]",
			"[password]EOL",
			"[data-value]EOF",
		))
	}

	var (
		storagePath = os.Args[2]
		dateKey     = []byte(os.Args[3])
	)

	sett := storage.NewSettings(&storage.SSettings{
		FPath:      storagePath,
		FWorkSize:  cWorkSize,
		FCipherKey: []byte(readLine("Password> ")),
	})
	stg, err := storage.NewCryptoStorage(sett)
	if err != nil {
		panic(err)
	}

	switch strings.ToUpper(os.Args[1]) {
	case "GET":
		data, err := stg.Get(dateKey)
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, bytes.NewBuffer(data))
	case "PUT":
		if _, err := stg.Get(dateKey); err == nil {
			panic("data-key already exist")
		}
		err := stg.Set(dateKey, readUntilEOF("Data> "))
		if err != nil {
			panic(err)
		}
		fmt.Println("ok")
	case "DEL":
		err := stg.Del(dateKey)
		if err != nil {
			panic(err)
		}
		fmt.Println("ok")
	default:
		panic("undefined option (get|put|del|new)")
	}
}

func readLine(s string) []byte {
	fmt.Print(s)
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}
	return line
}

func readUntilEOF(s string) []byte {
	fmt.Print(s)
	res, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return res
}
