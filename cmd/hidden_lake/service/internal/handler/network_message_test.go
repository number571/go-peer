package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
	anon_testutils "github.com/number571/go-peer/test/_data/anonymity"
)

func TestHandleMessageAPI(t *testing.T) {
	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[24])
	defer testAllFree(node, srv)

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[24]),
			&http.Client{Timeout: time.Minute},
		),
	)

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FWorkSize:    anon_testutils.TCWorkSize,
			FMessageSize: anon_testutils.TCMessageSize,
		}),
		privKey,
	)

	msg, err := client.EncryptPayload(
		privKey.GetPubKey(),
		payload.NewPayload(0, []byte("test")),
	)
	if err != nil {
		t.Error(err)
		return
	}

	panic(msg.ToBytes())
	// if err := hlsClient.HandleMessage(msg); err != nil {
	// 	t.Error(err)
	// 	return
	// }
}
