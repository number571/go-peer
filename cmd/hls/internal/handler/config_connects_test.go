package handler

import (
	"fmt"
	"testing"

	"github.com/number571/go-peer/cmd/hls/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
	"github.com/number571/go-peer/internal/testutils"
)

func TestHandleConnectsAPI(t *testing.T) {
	wcfg, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[6])
	defer testAllFree(node, srv)

	client := hls_client.NewClient(
		hls_client.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[6])),
	)

	connect := "test_connect4"
	testGetConnects(t, client, wcfg.Config())
	testAddConnect(t, client, connect)
	testDelConnect(t, client, connect)
}

func testGetConnects(t *testing.T, client hls_client.IClient, cfg config.IConfig) {
	connects, err := client.GetConnections()
	if err != nil {
		t.Error(err)
		return
	}

	if len(connects) != 3 {
		t.Errorf("length of connects != 3")
		return
	}

	for i := range connects {
		if connects[i] != cfg.Connections()[i] {
			t.Errorf("connections from config not equals with get")
			return
		}
	}
}

func testAddConnect(t *testing.T, client hls_client.IClient, connect string) {
	err := client.AddConnection(connect)
	if err != nil {
		t.Error(err)
		return
	}

	connects, err := client.GetConnections()
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
	err := client.DelConnection(connect)
	if err != nil {
		t.Error(err)
		return
	}

	connects, err := client.GetConnections()
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