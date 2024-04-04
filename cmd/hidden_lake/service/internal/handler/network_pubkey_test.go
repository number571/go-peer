package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandlePubKeyAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 8)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 8)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[8])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[8],
			&http.Client{Timeout: time.Minute},
		),
	)

	pubKey, err := client.GetPubKey(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if pubKey.ToString() != node.GetMessageQueue().GetClient().GetPubKey().ToString() {
		t.Error("public keys not equals")
		return
	}
}
