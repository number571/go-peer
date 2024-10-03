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

func TestHandleNetworkKeyAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 4)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 4)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[25])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[25],
			&http.Client{Timeout: time.Minute},
		),
	)

	testGetNetworkKey(t, client, "test")
}

func testGetNetworkKey(t *testing.T, client hls_client.IClient, networkKey string) {
	settings, err := client.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != networkKey {
		t.Error("got network key != networkKey")
		return
	}
}
