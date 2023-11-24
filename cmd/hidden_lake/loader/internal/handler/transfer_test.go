package handler

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"testing"
// 	"time"

// 	hll_client "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/client"
// 	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
// 	hls_app "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/app"
// 	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
// 	"github.com/number571/go-peer/pkg/client"
// 	"github.com/number571/go-peer/pkg/client/message"
// 	"github.com/number571/go-peer/pkg/crypto/asymmetric"
// 	net_message "github.com/number571/go-peer/pkg/network/message"
// 	"github.com/number571/go-peer/pkg/payload"
// 	testutils "github.com/number571/go-peer/test/_data"
// )

// const (
// 	tcTestData = "./test_data"
// 	tcNameHLT1 = tcTestData + "/hlt_1"
// 	tcNameHLT2 = tcTestData + "/hlt_2"
// )

// func TestTransferAPI(t *testing.T) {
// 	t.Parallel()

// 	netMsgSettings := net_message.NewSettings(&net_message.SSettings{
// 		FWorkSizeBits: testutils.TCWorkSize,
// 	})

// 	if err := copyWithPaste(tcNameHLT1, testutils.TgAddrs[42]); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer os.Remove(tcNameHLT1 + "/hlt.cfg")

// 	app1, err := hls_app.InitApp(tcNameHLT1)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err := app1.Run(); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	hltClient1 := hlt_client.NewClient(
// 		hlt_client.NewBuilder(),
// 		hlt_client.NewRequester(
// 			testutils.TgAddrs[42],
// 			&http.Client{Timeout: time.Minute / 2},
// 			netMsgSettings,
// 		),
// 	)

// 	msgSettings := message.NewSettings(
// 		&message.SSettings{
// 			FMessageSizeBytes: testutils.TCMessageSize,
// 		},
// 	)
// 	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
// 	client := client.NewClient(msgSettings, privKey)

// 	for i := 0; i < 5; i++ {
// 		encMsg, err := client.EncryptPayload(
// 			privKey.GetPubKey(),
// 			payload.NewPayload(
// 				uint64(i),
// 				[]byte("hello, world!"),
// 			),
// 		)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		err = hltClient1.PutMessage(
// 			net_message.NewMessage(
// 				netMsgSettings,
// 				payload.NewPayload(
// 					hls_settings.CNetworkMask,
// 					encMsg.ToBytes(),
// 				),
// 			),
// 		)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 	}

// 	if err := copyWithPaste(tcNameHLT2, testutils.TgAddrs[43]); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer os.Remove(tcNameHLT2 + "/hlt.cfg")

// 	app2, err := hls_app.InitApp(tcNameHLT2)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err := app2.Run(); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	service := testRunService(testutils.TgAddrs[44])
// 	defer service.Close()

// 	time.Sleep(100 * time.Millisecond)
// 	hllClient := hll_client.NewClient(
// 		hll_client.NewRequester(
// 			testutils.TgAddrs[44],
// 			&http.Client{Timeout: time.Second / 2},
// 		),
// 	)

// 	if err := hllClient.RunTransfer(); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	time.Sleep(5 * time.Second) // TODO
// }

// func copyWithPaste(pathTo, addr string) error {
// 	cfgDataFmt, err := os.ReadFile(tcTestData + "/hlt_copy.cfg")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(
// 		pathTo+"/hlt.cfg",
// 		[]byte(fmt.Sprintf(string(cfgDataFmt), testutils.TCWorkSize, addr)),
// 		0o644,
// 	)
// }
