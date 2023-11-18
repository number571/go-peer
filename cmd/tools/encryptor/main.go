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
	"github.com/number571/go-peer/pkg/crypto"
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
		mod     = strings.ToUpper(os.Args[1])
		param   crypto.IParameter
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
		FMessageSizeBytes: getMessageSize(param, dataValue),
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
		msg, err := message.LoadMessage(sett, dataValue)
		if err != nil {
			panic(err)
		}
		_, pld, err := c.DecryptMessage(msg)
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

func getMessageSize(param crypto.IParameter, dataValue []byte) uint64 {
	defaultSize := param.GetSize() << 1                // init size by key size
	return (defaultSize + uint64(len(dataValue))) << 1 // size with hex encode
}
