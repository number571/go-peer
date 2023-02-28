package handler

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleRequestAPI(t *testing.T) {
	_, node, srv := testAllCreate(tcPathConfig, tcPathDB, testutils.TgAddrs[9])
	defer testAllFree(node, srv)

	pushNode, pushSrv := testAllPushCreate()
	defer testAllPushFree(pushNode, pushSrv)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(fmt.Sprintf("http://%s", testutils.TgAddrs[9])),
	)

	node.GetNetworkNode().AddConnect(testutils.TgAddrs[11])
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	testBroadcast(t, client, pushNode.GetMessageQueue().GetClient().GetPubKey())
	testFetch(t, client, pushNode.GetMessageQueue().GetClient().GetPubKey())
}

func testBroadcast(t *testing.T, client hls_client.IClient, pubKey asymmetric.IPubKey) {
	err := client.BroadcastRequest(
		pubKey,
		request.NewRequest(http.MethodGet, tcServiceAddressInHLS, "/echo").
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

func testFetch(t *testing.T, client hls_client.IClient, pubKey asymmetric.IPubKey) {
	res, err := client.FetchRequest(
		pubKey,
		request.NewRequest(http.MethodGet, tcServiceAddressInHLS, "/echo").
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
	os.RemoveAll(tcPathConfig + "_push1")
	os.RemoveAll(tcPathDB + "_push1")

	pushNode := testNewPushNode(tcPathConfig+"_push1", tcPathDB+"_push1")
	if pushNode == nil {
		return nil, nil
	}

	pushSrv := testStartServerHTTP(testutils.TgAddrs[10])
	time.Sleep(200 * time.Millisecond)
	return pushNode, pushSrv
}

func testAllPushFree(node anonymity.INode, srv *http.Server) {
	defer func() {
		os.RemoveAll(tcPathConfig + "_push1")
		os.RemoveAll(tcPathDB + "_push1")
	}()
	types.StopAllCommands([]types.ICommand{
		node,
		node.GetNetworkNode(),
	})
	types.CloseAll([]types.ICloser{
		srv,
		node.GetKeyValueDB(),
	})
}

func testNewPushNode(cfgPath, dbPath string) anonymity.INode {
	node := testRunNewNode(dbPath, testutils.TgAddrs[11])
	rawCFG := &config.SConfig{
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[10],
		},
	}

	cfg, err := config.BuildConfig(cfgPath, rawCFG)
	if err != nil {
		return nil
	}

	node.HandleFunc(pkg_settings.CHeaderHLS, HandleServiceTCP(cfg))
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		if err := node.GetNetworkNode().Run(); err != nil {
			return
		}
	}()

	return node
}
