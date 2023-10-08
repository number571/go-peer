package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandlePubKeyAPI(t *testing.T) {
	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[8])
	defer testAllFree(node, srv)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[8]),
			&http.Client{Timeout: time.Minute},
		),
	)

	pubKey, ephPubKey, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.ToString() != node.GetMessageQueue().GetClient().GetPubKey().ToString() {
		t.Error("public keys not equals")
		return
	}

	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc2PrivKey1024)
	if err := client.SetPrivKey(privKey, ephPubKey); err != nil {
		t.Error("failed update private key")
		return
	}

	newPubKey, _, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.GetAddress().ToString() == newPubKey.GetAddress().ToString() {
		t.Error("public keys are equals")
		return
	}

	if err := client.ResetPrivKey(); err != nil {
		t.Error(err)
		return
	}

	gotInitPubKey, _, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}

	// tgInitPrivKey = testutils.Tc1PrivKey1024
	if gotInitPubKey.GetAddress().ToString() != tgInitPrivKey.GetPubKey().GetAddress().ToString() {
		t.Error("init state of private key is incorrect after reset")
		return
	}
}
