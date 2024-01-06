package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/closer"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestHandleRequestAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 7)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 7)

	_, node, ctx, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[9])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	pushNode, pushCancel, pushSrv := testAllPushCreate(pathCfg, pathDB)
	defer testAllPushFree(pushNode, pushCancel, pushSrv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[9]),
			&http.Client{Timeout: time.Minute},
		),
	)

	node.GetNetworkNode().AddConnection(ctx, testutils.TgAddrs[11])
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024).GetPubKey())

	testBroadcast(t, client)
	testFetch(t, client)
}

func testBroadcast(t *testing.T, client hls_client.IClient) {
	err := client.BroadcastRequest(
		"test_recvr",
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

func testFetch(t *testing.T, client hls_client.IClient) {
	res, err := client.FetchRequest(
		"test_recvr",
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

	body := res.GetBody()
	if string(body) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		t.Errorf("result does not match; got '%s'", string(body))
		return
	}
}

func testAllPushCreate(pathCfg, pathDB string) (anonymity.INode, context.CancelFunc, *http.Server) {
	os.RemoveAll(pathCfg + "_push1")
	os.RemoveAll(pathDB + "_push1")

	pushNode, cancel := testNewPushNode(pathCfg+"_push1", pathDB+"_push1")
	if pushNode == nil {
		return nil, cancel, nil
	}

	pushSrv := testStartServerHTTP(testutils.TgAddrs[10])
	time.Sleep(200 * time.Millisecond)
	return pushNode, cancel, pushSrv
}

func testAllPushFree(node anonymity.INode, cancel context.CancelFunc, srv *http.Server, pathCfg, pathDB string) {
	defer func() {
		os.RemoveAll(pathCfg + "_push1")
		os.RemoveAll(pathDB + "_push1")
	}()
	cancel()
	closer.CloseAll([]types.ICloser{
		srv,
		node.GetDBWrapper(),
		node.GetNetworkNode(),
	})
}

func testNewPushNode(cfgPath, dbPath string) (anonymity.INode, context.CancelFunc) {
	node, ctx, cancel := testRunNewNode(dbPath, testutils.TgAddrs[11])
	rawCFG := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FKeySizeBits:      testutils.TcKeySize,
			FQueuePeriodMS:    testutils.TCQueuePeriod,
		},
		FServices: map[string]*config.SService{
			tcServiceAddressInHLS: {FHost: testutils.TgAddrs[10]},
		},
	}

	cfg, err := config.BuildConfig(cfgPath, rawCFG)
	if err != nil {
		return nil, cancel
	}

	node.HandleFunc(
		hls_settings.CServiceMask,
		HandleServiceTCP(
			&sync.Mutex{},
			config.NewWrapper(cfg),
			logger.NewLogger(
				logger.NewSettings(&logger.SSettings{}),
				func(_ logger.ILogArg) string { return "" },
			),
		),
	)
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024).GetPubKey())

	go func() { _ = node.GetNetworkNode().Listen(ctx) }()

	return node, cancel
}
