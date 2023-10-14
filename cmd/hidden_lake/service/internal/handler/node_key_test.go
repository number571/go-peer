package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
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

	pubKey, err := client.GetPubKey()
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.ToString() != node.GetMessageQueue().GetClient().GetPubKey().ToString() {
		t.Error("public keys not equals")
		return
	}
}
