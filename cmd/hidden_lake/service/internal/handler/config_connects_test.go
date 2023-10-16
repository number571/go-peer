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

func TestHandleConnectsAPI(t *testing.T) {
	wcfg, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[6])
	defer testAllFree(node, srv)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[6]),
			&http.Client{Timeout: time.Minute},
		),
	)

	connect := "test_connect4"
	testGetConnects(t, client, wcfg.GetConfig())
	testAddConnect(t, client, connect)
	testDelConnect(t, client, connect)
}

func testGetConnects(t *testing.T, client hls_client.IClient, cfg config.IConfig) {
	connects, err := client.GetConnections(false)
	if err != nil {
		t.Error(err)
		return
	}

	if len(connects) != 3 {
		t.Error("length of connects != 3")
		return
	}

	for i := range connects {
		if connects[i] != cfg.GetConnections()[i] {
			t.Error("connections from config not equals with get")
			return
		}
	}
}

func testAddConnect(t *testing.T, client hls_client.IClient, connect string) {
	err := client.AddConnection(false, connect)
	if err != nil {
		t.Error(err)
		return
	}

	connects, err := client.GetConnections(false)
	if err != nil {
		t.Error(err)
		return
	}

	for _, conn := range connects {
		if conn == connect {
			return
		}
	}
	t.Errorf("undefined connection key by '%s'", connect)
}

func testDelConnect(t *testing.T, client hls_client.IClient, connect string) {
	err := client.DelConnection(false, connect)
	if err != nil {
		t.Error(err)
		return
	}

	connects, err := client.GetConnections(false)
	if err != nil {
		t.Error(err)
		return
	}

	for _, conn := range connects {
		if conn == connect {
			t.Errorf("deleted connect exists for '%s'", connect)
			return
		}
	}
}
