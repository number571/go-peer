package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	"github.com/number571/go-peer/internal/settings"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	pldHead = 0x1
)

func main() {
	cfg := &settings.SConfigSettings{
		FSettings: settings.SConfigSettingsBlock{
			FWorkSizeBits:     20,
			FMessageSizeBytes: (8 << 10),
		},
	}

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://localhost:9582",
			&http.Client{Timeout: time.Minute},
			message.NewSettings(&message.SSettings{
				FWorkSizeBits:     cfg.GetWorkSizeBits(),
				FMessageSizeBytes: cfg.GetMessageSizeBytes(),
			}),
		),
	)

	readPrivKey, err := filesystem.OpenFile("priv.key").Read()
	if err != nil {
		panic(err)
	}

	privKey := asymmetric.LoadRSAPrivKey(string(readPrivKey))
	client := hls_settings.InitClient(cfg, privKey)

	if len(os.Args) < 2 {
		panic("len os.Args < 2")
	}

	switch os.Args[1] {
	case "w", "write":
		if len(os.Args) != 3 {
			panic("len os.Args != 3")
		}

		msg, err := client.EncryptPayload(
			privKey.GetPubKey(),
			payload.NewPayload(pldHead, []byte(os.Args[2])),
		)
		if err != nil {
			panic(err)
		}

		if err := hltClient.PutMessage(msg); err != nil {
			panic(err)
		}

		fmt.Printf("%x\n", msg.GetBody().GetHash())
	case "r", "read":
		if len(os.Args) != 3 {
			panic("len os.Args != 3")
		}

		msg, err := hltClient.GetMessage(os.Args[2])
		if err != nil {
			panic(err)
		}

		pubKey, pld, err := client.DecryptMessage(msg)
		if err != nil {
			panic(err)
		}

		if pld.GetHead() != pldHead {
			panic("payload head != constant head")
		}

		if pubKey.GetAddress().ToString() != client.GetPubKey().GetAddress().ToString() {
			panic("public key is incorrect")
		}

		fmt.Println(string(pld.GetBody()))
	case "h", "hashes":
		hashes, err := hltClient.GetHashes()
		if err != nil {
			panic(err)
		}

		for i, hash := range hashes {
			fmt.Printf("%d. %s\n", i, hash)
		}
	}
}
