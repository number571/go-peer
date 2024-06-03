package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cPldHead     = 0x1
	cKeySize     = 4096
	cMsgSize     = (8 << 10)
	cWrkSize     = 22
	cPrivKeyPath = "../_keys/priv.key"
)

var gAddrHLTs = [][2]string{
	{"94.103.91.81:9582", "8Jkl93Mdk93md1bz"},
	{"195.133.1.126:9582", "kf92j74Nof92n9F4"},
	{"193.233.18.245:9582", "oi4r9NW9Le7fKF9d"},
	{"185.43.4.253:9582", "j2BR39JfDf7Bajx3"},
}

func main() {
	ctx := context.Background()
	randString := random.NewCSPRNG().GetString(16)

	sett := message.NewSettings(&message.SSettings{
		FMessageSizeBytes: cMsgSize,
		FKeySizeBits:      cKeySize,
	})

	readPrivKey, err := os.ReadFile(cPrivKeyPath)
	if err != nil {
		panic(err)
	}

	privKey := asymmetric.LoadRSAPrivKey(string(readPrivKey))
	client := client.NewClient(sett, privKey)

	for i, addrHLT := range gAddrHLTs {
		netSett := net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: cWrkSize,
			FNetworkKey:   addrHLT[1],
		})

		hltClient := hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				"http://"+addrHLT[0],
				&http.Client{Timeout: 5 * time.Second},
				netSett,
			),
		)

		msg, err := client.EncryptMessage(
			privKey.GetPubKey(),
			payload.NewPayload64(cPldHead, []byte(randString)).ToBytes(),
		)
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i+1, addrHLT[0], err)
			continue
		}

		netMsg := net_message.NewMessage(
			netSett,
			payload.NewPayload32(hls_settings.CNetworkMask, msg),
		)

		if err := hltClient.PutMessage(ctx, netMsg); err != nil {
			fmt.Printf("%d. %s: %s\n", i+1, addrHLT[0], err)
			continue
		}

		gotNetMsg, err := hltClient.GetMessage(ctx, encoding.HexEncode(netMsg.GetHash()))
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i+1, addrHLT[0], err)
			continue
		}

		if !bytes.Equal(netMsg.ToBytes(), gotNetMsg.ToBytes()) {
			fmt.Printf("%d. %s: !bytes.Equal(netMsg.ToBytes(), gotNetMsg.ToBytes())\n", i+1, addrHLT[0])
			continue
		}

		fmt.Printf("%d. HLT server '%s' is working properly;\n", i+1, addrHLT[0])
	}
}
