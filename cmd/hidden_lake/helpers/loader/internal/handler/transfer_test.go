package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	hll_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	hls_app "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/app"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcMessageSize = (8 << 10)
)

const (
	tcTestData = "./test_data"
	tcNameHLT1 = tcTestData + "/hlt_1"
	tcNameHLT2 = tcTestData + "/hlt_2"
)

func testCreateHLT(netMsgSettings net_message.ISettings, path, addr string) (types.IRunner, context.CancelFunc, hlt_client.IClient, error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := copyWithPaste(path, addr); err != nil {
		return nil, cancel, nil, err
	}

	app1, err := hls_app.InitApp(path)
	if err != nil {
		return nil, cancel, nil, err
	}

	go func() { _ = app1.Run(ctx) }()
	time.Sleep(100 * time.Millisecond)

	hltClient1 := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute / 2},
			netMsgSettings,
		),
	)

	return app1, cancel, hltClient1, nil
}

func testInitTransfer() {
	os.RemoveAll(tcNameHLT1)
	os.RemoveAll(tcNameHLT2)

	os.Mkdir(tcNameHLT1, 0o777)
	os.Mkdir(tcNameHLT2, 0o777)

}

func TestHandleTransferAPI(t *testing.T) {
	t.Parallel()

	testInitTransfer()
	defer func() {
		os.RemoveAll(tcNameHLT1)
		os.RemoveAll(tcNameHLT2)
	}()

	// INIT SERVICES

	netMsgSettings := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: testutils.TCWorkSize,
		FNetworkKey:   testutils.TCNetworkKey,
	})

	_, cancel1, hltClient1, err := testCreateHLT(netMsgSettings, tcNameHLT1, tgProducer)
	if err != nil {
		t.Error(err)
		return
	}
	defer cancel1()

	_, cancel2, hltClient2, err := testCreateHLT(netMsgSettings, tcNameHLT2, tgConsumer)
	if err != nil {
		t.Error(err)
		return
	}
	defer cancel2()

	service := testRunService(tgTService)
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+tgTService,
			&http.Client{Timeout: time.Second / 2},
		),
	)

	// PUSH MESSAGES

	msgSettings := message.NewSettings(
		&message.SSettings{
			FMessageSizeBytes: tcMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		},
	)
	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	client := client.NewClient(msgSettings, privKey)

	for i := 0; i < 5; i++ {
		encMsg, err := client.EncryptPayload(
			privKey.GetPubKey(),
			payload.NewPayload(
				uint64(i),
				[]byte("hello, world!"),
			),
		)
		if err != nil {
			t.Error(err)
			return
		}
		netMsg := net_message.NewMessage(
			netMsgSettings,
			payload.NewPayload(
				hls_settings.CNetworkMask,
				encMsg.ToBytes(),
			),
			1,
		)
		err = hltClient1.PutMessage(netMsg)
		if err != nil {
			t.Error(err)
			return
		}
	}

	// TRANSFER MESSAGES

	if err := hllClient.RunTransfer(); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)

	// LOAD MESSAGES

	for i := uint64(0); ; i++ {
		h, err := hltClient2.GetHash(i)
		if err != nil {
			if i != 5 {
				t.Error(i, err)
			}
			return
		}

		netMsg, err := hltClient2.GetMessage(h)
		if err != nil {
			t.Error(err)
			return
		}

		msg, err := message.LoadMessage(msgSettings, netMsg.GetPayload().GetBody())
		if err != nil {
			t.Error(err)
			return
		}

		pubKey, pld, err := client.DecryptMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}

		if pld.GetHead() != uint64(i) {
			t.Error("got bad index")
			return
		}

		if pubKey.GetHasher().ToString() != client.GetPubKey().GetHasher().ToString() {
			t.Error("got bad public key")
			return
		}
	}
}

func copyWithPaste(pathTo, addr string) error {
	cfgDataFmt, err := os.ReadFile(tcTestData + "/hlt_copy.yml")
	if err != nil {
		return err
	}
	return os.WriteFile(
		pathTo+"/hlt.yml",
		[]byte(fmt.Sprintf(string(cfgDataFmt), testutils.TCWorkSize, testutils.TCNetworkKey, addr)),
		0o644,
	)
}
