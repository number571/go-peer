package handler

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/hlc"
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/settings/testutils"
)

func TestHandlePushAPI(t *testing.T) {
	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[9])
	defer testAllFree(node, srv)

	pushNode, pushSrv := testAllPushCreate()
	defer testAllPushFree(pushNode, pushSrv)

	client := hlc.NewClient(
		hlc.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[9])),
	)

	node.Network().Connect(testutils.TgAddrs[11])
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	testBroadcast(t, client, pushNode.Queue().Client().PubKey())
	testRequest(t, client, pushNode.Queue().Client().PubKey())
}

func testBroadcast(t *testing.T, client hlc.IClient, pubKey asymmetric.IPubKey) {
	err := client.Broadcast(
		pubKey,
		hls_network.NewRequest("GET", tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)),
	)
	if err != nil {
		t.Error(err)
		return
	}
}

func testRequest(t *testing.T, client hlc.IClient, pubKey asymmetric.IPubKey) {
	res, err := client.Request(
		pubKey,
		hls_network.NewRequest("GET", tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if string(res) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		t.Errorf("result does not match; get '%s'", string(res))
		return
	}
}

func testAllPushCreate() (anonymity.INode, *http.Server) {
	pushNode := testNewPushNode(tcPathConfig+"_push1", tcPathDB+"_push1")
	if pushNode == nil {
		return nil, nil
	}
	pushSrv := testStartServerHTTP()
	time.Sleep(200 * time.Millisecond)
	return pushNode, pushSrv
}

func testAllPushFree(node anonymity.INode, srv *http.Server) {
	defer func() {
		os.RemoveAll(tcPathConfig + "_push1")
		os.RemoveAll(tcPathDB + "_push1")
		closer.CloseAll([]modules.ICloser{
			node,
			srv,
			node.KeyValueDB(),
			node.Network(),
		})
	}()
}

func testNewPushNode(cfgPath, dbPath string) anonymity.INode {
	node := testRunNewNode(dbPath)
	rawCFG := &config.SConfig{
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[10],
		},
	}

	cfg, err := config.NewConfig(cfgPath, rawCFG)
	if err != nil {
		return nil
	}

	node.Handle(hls_settings.CHeaderHLS, HandleServiceTCP(cfg))
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		err := node.Network().Listen(testutils.TgAddrs[11])
		if err != nil {
			return
		}
	}()

	return node
}

func testStartServerHTTP() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testutils.TestEchoPage)

	srv := &http.Server{
		Addr:    testutils.TgAddrs[10],
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}
