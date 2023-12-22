package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	hll_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	hltHost1 = "localhost:7582"
	hltHost2 = "localhost:7582"
	hllHost  = "localhost:6561"
)

const (
	messageSize = (8 << 10) // 8KiB
	networkKey  = "some-network-key"
	workSize    = 10
	keySize     = 1024
)

const (
	messageCount = 64
)

var (
	privKey      asymmetric.IPrivKey
	pushedHashes = make([][]byte, 0, messageCount)
)

func init() {
	readPrivKey, err := os.ReadFile("priv.key")
	if err != nil {
		panic(err)
	}
	privKey = asymmetric.LoadRSAPrivKey(string(readPrivKey))
}

func main() {
	netMsgSettings := net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   networkKey,
		FWorkSizeBits: workSize,
	})

	msgSettings := message.NewSettings(&message.SSettings{
		FMessageSizeBytes: messageSize,
		FKeySizeBits:      keySize,
	})

	if err := pushMessages(netMsgSettings, msgSettings); err != nil {
		panic(err)
	}

	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+hllHost,
			&http.Client{Timeout: time.Minute / 2},
		),
	)
	if err := hllClient.RunTransfer(); err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	if err := checkMessages(netMsgSettings, msgSettings); err != nil {
		panic(err)
	}

	fmt.Println("messages have been successfully transported")
}

func pushMessages(netMsgSettings net_message.ISettings, msgSettings message.ISettings) error {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+hltHost1,
			&http.Client{Timeout: time.Minute / 2},
			netMsgSettings,
		),
	)

	client := client.NewClient(msgSettings, privKey)

	for i := 0; i < messageCount; i++ {
		msg, err := client.EncryptPayload(
			client.GetPubKey(), // self encrypt
			payload.NewPayload(uint64(i), []byte("hello, world!")),
		)
		if err != nil {
			return err
		}

		netMsg := net_message.NewMessage(
			netMsgSettings,
			payload.NewPayload(hls_settings.CNetworkMask, msg.ToBytes()),
		)
		if err := hltClient.PutMessage(netMsg); err != nil {
			return err
		}

		pushedHashes = append(pushedHashes, netMsg.GetHash())
	}

	return nil
}

func checkMessages(netMsgSettings net_message.ISettings, msgSettings message.ISettings) error {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+hltHost2,
			&http.Client{Timeout: time.Minute / 2},
			netMsgSettings,
		),
	)

	client := client.NewClient(msgSettings, privKey)

	hashes := make([]string, 0, messageCount)
	for i := uint64(0); ; i++ {
		hash, err := hltClient.GetHash(i)
		if err != nil {
			break
		}
		hashes = append(hashes, hash)
	}

	for _, ph := range pushedHashes {
		if !hashIsExist(ph, hashes) {
			return errors.New("hash not found")
		}
	}

	for _, h := range hashes {
		netMsg, err := hltClient.GetMessage(h)
		if err != nil {
			return err
		}

		if netMsg.GetPayload().GetHead() != hls_settings.CNetworkMask {
			return errors.New("network mask is invalid")
		}

		encMsg, err := message.LoadMessage(msgSettings, netMsg.GetPayload().GetBody())
		if err != nil {
			return err
		}

		pubKey, pld, err := client.DecryptMessage(encMsg)
		if err != nil {
			return err
		}

		if pubKey.GetAddress().ToString() != client.GetPubKey().GetAddress().ToString() {
			return errors.New("got invalid public key")
		}

		if pld.GetHead() > messageCount {
			return errors.New("got invalid head value")
		}

		if string(pld.GetBody()) != "hello, world!" {
			return errors.New("got invalid body value")
		}
	}

	return nil
}
func hashIsExist(hash []byte, listHashes []string) bool {
	strHash := encoding.HexEncode(hash)
	for _, h := range listHashes {
		if strHash == h {
			return true
		}
	}
	return false
}
