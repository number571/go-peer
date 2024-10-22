package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/number571/go-peer/pkg/crypto/keybuilder"
)

const (
	cKeySize = 33 // bytes
)

func main() {
	saltParam := flag.String("salt", "_salt_", "default salt value")
	workParam := flag.Uint("work", 24, "default work value")
	flag.Parse()

	keyBuilder := keybuilder.NewKeyBuilder(1<<(*workParam), []byte(*saltParam))
	gotPassword := keyBuilder.Build(readUntilEOL(), cKeySize)

	fmt.Println(base64.URLEncoding.EncodeToString(gotPassword))
}

func readUntilEOL() string {
	var (
		p = make([]byte, 0, 256)
		b = make([]byte, 1)
	)

	if runtime.GOOS == "windows" {
		res, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			panic(err)
		}
		return string(res)
	}

	if err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("stty", "-F", "/dev/tty", "-echo").Run(); err != nil {
		panic(err)
	}
	defer func() { _ = exec.Command("stty", "-F", "/dev/tty", "echo").Run() }()

	for {
		if _, err := os.Stdin.Read(b); err != nil {
			panic(err)
		}
		if b[0] == '\n' { // <enter>
			fmt.Println()
			break
		}
		if b[0] == 127 { // <backspace>
			if len(p) == 0 {
				continue
			}
			fmt.Print("\r")
			for i := 0; i < len(p); i++ {
				fmt.Print(" ")
			}
			fmt.Print("\r")
			for i := 0; i < len(p)-1; i++ {
				fmt.Print("*")
			}
			p = p[:len(p)-1]
			continue
		}
		fmt.Print("*")
		p = append(p, b[0])
	}

	return string(p)
}
