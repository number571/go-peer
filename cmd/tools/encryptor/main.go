package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Sprintf(
			"usage:\n\t%s\nstdin:\n\t%s\n\n",
			"./main (e|d) [pubkey-file|privkey-file]",
			"[data-value]EOF",
		))
	}

	var (
		mode    = strings.ToUpper(os.Args[1])
		pubKey  asymmetric.IPubKey
		privKey asymmetric.IPrivKey
		keySize uint64
	)

	keyBytes, err := os.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	switch mode {
	case "E":
		pubKey = asymmetric.LoadRSAPubKey(string(keyBytes))
		if pubKey == nil {
			panic("incorrect public key")
		}
		keySize = pubKey.GetSize()
	case "D":
		privKey = asymmetric.LoadRSAPrivKey(string(keyBytes))
		if privKey == nil {
			panic("incorrect private key")
		}
		keySize = privKey.GetSize()
	default:
		panic("unkown mode")
	}

	// TODO: dynamic size of message
	dataValue := string(readUntilEOF())
	sett := message.NewSettings(&message.SSettings{
		FMessageSizeBytes: (8 << 10),
		FKeySizeBits:      keySize,
	})

	switch mode {
	case "E":
		c := client.NewClient(sett, asymmetric.NewRSAPrivKey(pubKey.GetSize()))
		msg, err := c.EncryptPayload(pubKey, payload.NewPayload(0x1, []byte(dataValue)))
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.ToString())
	case "D":
		c := client.NewClient(sett, privKey)
		msg, err := message.LoadMessage(sett, dataValue)
		if err != nil {
			panic(err)
		}
		_, pld, err := c.DecryptMessage(msg)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(os.Stdout, bytes.NewBuffer(pld.GetBody())); err != nil {
			panic(err)
		}
	}
}

func readUntilEOF() []byte {
	res, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return res
}
