package handler

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hls/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func TestHandleOnlineAPI(t *testing.T) {
	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[12])
	defer testAllFree(node, srv)

	pushNode := testAllOnlineCreate()
	defer testAllOnlineFree(pushNode)

	client := hls_client.NewClient(
		hls_client.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[12])),
	)

	node.Network().Connect(testutils.TgAddrs[13])
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	testGetOnlines(t, client, node)
	testDelOnline(t, client, testutils.TgAddrs[13])
}

func testGetOnlines(t *testing.T, client hls_client.IClient, node anonymity.INode) {
	onlines, err := client.GetOnlines()
	if err != nil {
		t.Error(err)
		return
	}

	if len(onlines) != 1 {
		t.Errorf("length of onlines != 1")
		return
	}

	if _, ok := node.Network().Connections()[onlines[0]]; !ok {
		fmt.Println(node.Network().Connections(), onlines[0])
		t.Errorf("online address is invalid")
		return
	}
}

func testDelOnline(t *testing.T, client hls_client.IClient, addr string) {
	err := client.DelConnection(addr)
	if err != nil {
		t.Error(err)
		return
	}

	onlines, err := client.GetOnlines()
	if err != nil {
		t.Error(err)
		return
	}

	if len(onlines) != 0 {
		t.Errorf("length of onlines != 0")
		return
	}
}

func testAllOnlineCreate() anonymity.INode {
	pushNode := testOnlinePushNode(tcPathConfig+"_push2", tcPathDB+"_push2")
	if pushNode == nil {
		return nil
	}
	time.Sleep(200 * time.Millisecond)
	return pushNode
}

func testAllOnlineFree(node anonymity.INode) {
	defer func() {
		os.RemoveAll(tcPathConfig + "_push2")
		os.RemoveAll(tcPathDB + "_push2")
		closer.CloseAll([]modules.ICloser{
			node,
			node.KeyValueDB(),
			node.Network(),
		})
	}()
}

func testOnlinePushNode(cfgPath, dbPath string) anonymity.INode {
	node := testRunNewNode(dbPath)

	cfg, err := config.NewConfig(cfgPath, &config.SConfig{})
	if err != nil {
		return nil
	}

	node.Handle(hls_settings.CHeaderHLS, HandleServiceTCP(cfg))
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		err := node.Network().Listen(testutils.TgAddrs[13])
		if err != nil {
			return
		}
	}()

	return node
}
