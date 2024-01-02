package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	hle_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

const (
	cLocalAddressHLE = "localhost:8551"
)

func main() {
	if len(os.Args) != 3 {
		panic("len(os.Args) != 3")
	}

	netSett := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: 20,
		FNetworkKey:   "j2BR39JfDf7Bajx3",
	})

	hleClient := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+cLocalAddressHLE,
			&http.Client{Timeout: time.Minute},
			netSett,
		),
	)

	switch os.Args[1] {
	case "e", "encrypt":
		readPubKey, err := os.ReadFile("pub.key")
		if err != nil {
			panic(err)
		}

		pubKey := asymmetric.LoadRSAPubKey(string(readPubKey))
		netMsg, err := hleClient.EncryptMessage(pubKey, []byte(os.Args[2]))
		if err != nil {
			panic(err)
		}

		fmt.Println(netMsg.ToString())
	case "d", "decrypt":
		netMsg, err := net_message.LoadMessage(netSett, string(os.Args[2]))
		if err != nil {
			panic(err)
		}

		_, data, err := hleClient.DecryptMessage(netMsg)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))
	default:
		panic("unknown mode")
	}
}
