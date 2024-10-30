package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
)

func main() {
	seed := flag.Bool("seed", false, "set seed private key")
	flag.Parse()

	seedBytes := random.NewRandom().GetBytes(asymmetric.CKeySeedSize)
	if *seed {
		seedBytes = encoding.HexDecode(readUntilEOL())
		if len(seedBytes) != asymmetric.CKeySeedSize {
			panic("len(seedBytes) != asymmetric.CKeySeedSize")
		}
	}

	priv := asymmetric.NewPrivKeyFromSeed(seedBytes)
	if err := os.WriteFile("seed.key", []byte(encoding.HexEncode(seedBytes)), 0o600); err != nil {
		panic(err)
	}
	if err := os.WriteFile("priv.key", []byte(priv.ToString()), 0o600); err != nil {
		panic(err)
	}
	if err := os.WriteFile("pub.key", []byte(priv.GetPubKey().ToString()), 0o600); err != nil {
		panic(err)
	}
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
