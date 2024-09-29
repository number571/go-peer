package handler

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	hle_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleEncryptDecryptAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[48])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[48],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	// same private key in the HLE
	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024)
	pubKey := privKey.GetPubKey()

	data := []byte("hello, world!")

	netMsg, err := hleClient.EncryptMessage(context.Background(), pubKey, payload.NewPayload64(1, data))
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, getPld, err := hleClient.DecryptMessage(context.Background(), netMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(pubKey.ToBytes(), gotPubKey.ToBytes()) {
		t.Error("got invalid public key")
		return
	}

	if !bytes.Equal(getPld.GetBody(), data) {
		t.Error("got invalid data")
		return
	}
}
