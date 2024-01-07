package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cPldHead = 0x1
)

const (
	cLocalAddressHLT = "localhost:8582"
	cProdAddressHLT  = "185.43.4.253:9582" // 1x3.1GHz, 2.0GB RAM, 300GB HDD
)

func main() {
	sett := message.NewSettings(&message.SSettings{
		FMessageSizeBytes: (8 << 10),
		FKeySizeBits:      4096,
	})

	netSett := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: 22,
		FNetworkKey:   "j2BR39JfDf7Bajx3",
	})

	readPrivKey, err := os.ReadFile("priv.key")
	if err != nil {
		panic(err)
	}

	privKey := asymmetric.LoadRSAPrivKey(string(readPrivKey))
	client := client.NewClient(sett, privKey)

	if len(os.Args) < 2 {
		panic("len os.Args < 2")
	}

	args := os.Args[1:]
	addr := cLocalAddressHLT

	if args[0] == "prod" {
		args = os.Args[2:]
		addr = cProdAddressHLT
		if len(args) == 0 {
			panic("len os.Args < 2")
		}
	}

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute},
			netSett,
		),
	)

	switch args[0] {
	case "w", "write":
		if len(args) != 2 {
			panic("len write.args != 2")
		}

		msg, err := client.EncryptPayload(
			privKey.GetPubKey(),
			payload.NewPayload(cPldHead, []byte(args[1])),
		)
		if err != nil {
			panic(err)
		}

		netMsg := net_message.NewMessage(
			netSett,
			payload.NewPayload(hls_settings.CNetworkMask, msg.ToBytes()),
			1,
		)

		if err := hltClient.PutMessage(netMsg); err != nil {
			panic(err)
		}

		fmt.Printf("%x\n", netMsg.GetHash())
	case "r", "read":
		if len(args) != 2 {
			panic("len read.args != 2")
		}

		netMsg, err := hltClient.GetMessage(args[1])
		if err != nil {
			panic(err)
		}

		msg, err := message.LoadMessage(client.GetSettings(), netMsg.GetPayload().GetBody())
		if err != nil {
			panic("load message is nil")
		}

		pubKey, pld, err := client.DecryptMessage(msg)
		if err != nil {
			panic(err)
		}

		if pld.GetHead() != cPldHead {
			panic("payload head != constant head")
		}

		if pubKey.GetHasher().ToString() != client.GetPubKey().GetHasher().ToString() {
			panic("public key is incorrect")
		}

		fmt.Println(string(pld.GetBody()))
	case "h", "hashes":
		for i := uint64(0); ; i++ {
			hash, err := hltClient.GetHash(i)
			if err != nil {
				return
			}
			fmt.Printf("[%d] %s\n", i+1, hash)
		}
	default:
		panic("unknown mode")
	}
}
