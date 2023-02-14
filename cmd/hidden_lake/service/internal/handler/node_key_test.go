package handler

import (
	"fmt"
	"testing"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandlePubKeyAPI(t *testing.T) {
	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[8])
	defer testAllFree(node, srv)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[8])),
	)

	pubKey, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.String() != node.Queue().Client().PubKey().String() {
		t.Error("public keys not equals")
		return
	}

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey2)
	if err := client.SetPrivKey(privKey); err != nil {
		t.Error("failed update private key")
		return
	}

	newPubKey, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.Address().String() == newPubKey.Address().String() {
		t.Error("public keys are equals")
		return
	}
}
