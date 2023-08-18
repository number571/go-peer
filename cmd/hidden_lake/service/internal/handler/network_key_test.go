package handler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleNetworkKeyAPI(t *testing.T) {
	wcfg, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[26])
	defer testAllFree(node, srv)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[26]),
			&http.Client{Timeout: time.Minute},
		),
	)

	networkKey := "test_network_key"
	testSetNetworkKey(t, client, networkKey)
	testGetNetworkKey(t, client, wcfg.GetConfig(), networkKey)
}

func testGetNetworkKey(t *testing.T, client hls_client.IClient, cfg config.IConfig, networkKey string) {
	gotNetworkKey, err := client.GetNetworkKey()
	if err != nil {
		t.Error(err)
		return
	}

	if gotNetworkKey != networkKey {
		t.Error("got network key != networkKey")
		return
	}
}

func testSetNetworkKey(t *testing.T, client hls_client.IClient, networkKey string) {
	err := client.SetNetworkKey(networkKey)
	if err != nil {
		t.Error(err)
		return
	}
}