package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

const (
	cWorkSize = 20 // bits
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
		mod     = strings.ToUpper(os.Args[1])
		param   types.IParameter
		pubKey  asymmetric.IPubKey
		privKey asymmetric.IPrivKey
	)

	keyBytes, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	switch mod {
	case "E":
		pubKey = asymmetric.LoadRSAPubKey(string(keyBytes))
		if pubKey == nil {
			panic("incorrect public key")
		}
		param = pubKey
	case "D":
		privKey = asymmetric.LoadRSAPrivKey(string(keyBytes))
		if privKey == nil {
			panic("incorrect private key")
		}
		param = privKey
	}

	dataValue := readUntilEOF()
	sett := message.NewSettings(&message.SSettings{
		FWorkSize:    cWorkSize,
		FMessageSize: getMessageSize(param, dataValue),
	})

	switch mod {
	case "E":
		c := client.NewClient(sett, asymmetric.NewRSAPrivKey(pubKey.GetSize()))
		msg, err := c.EncryptPayload(pubKey, payload.NewPayload(0x1, dataValue))
		if err != nil {
			panic(err)
		}
		fmt.Println(string(msg.ToBytes()))
	case "D":
		c := client.NewClient(sett, privKey)
		_, pld, err := c.DecryptMessage(message.LoadMessage(sett, dataValue))
		if err != nil {
			panic(err)
		}
		io.Copy(os.Stdout, bytes.NewBuffer(pld.GetBody()))
	}
}

func readUntilEOF() []byte {
	res, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return res
}

func getMessageSize(param types.IParameter, dataValue []byte) uint64 {
	defaultSize := (4 << 10) + (param.GetSize() / 8)   // in bytes with sign+pubKey
	return (defaultSize + uint64(len(dataValue))) << 1 // size with hex encode
}
