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
)

func TestHandleMessageAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 5)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 5)

	_, node, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[24])
	defer testAllFree(node, srv, pathCfg, pathDB)

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[24]),
			&http.Client{Timeout: time.Minute},
		),
	)

	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	pubKey := asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0])

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
		}),
		privKey,
	)
	msg, err := client.EncryptPayload(
		pubKey,
		payload.NewPayload(0, []byte("hello")),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := hlsClient.HandleMessage(msg); err != nil {
		t.Error(err)
		return
	}
}
