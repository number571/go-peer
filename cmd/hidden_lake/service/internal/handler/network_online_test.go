package handler

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestHandleOnlineAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 6)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 6)

	_, node, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[12])
	defer testAllFree(node, srv, pathCfg, pathDB)

	pushNode := testAllOnlineCreate(pathCfg, pathDB)
	defer testAllOnlineFree(pushNode, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", testutils.TgAddrs[12]),
			&http.Client{Timeout: time.Minute},
		),
	)

	node.GetNetworkNode().AddConnection(testutils.TgAddrs[13])
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024).GetPubKey())

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
		t.Error("length of onlines != 1")
		return
	}

	if _, ok := node.GetNetworkNode().GetConnections()[onlines[0]]; !ok {
		t.Error("online address is invalid")
		return
	}
}

func testDelOnline(t *testing.T, client hls_client.IClient, addr string) {
	err := client.DelConnection(false, addr)
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
		t.Error("length of onlines != 0")
		return
	}
}

func testAllOnlineCreate(pathCfg, pathDB string) anonymity.INode {
	os.RemoveAll(pathCfg + "_push2")
	os.RemoveAll(pathDB + "_push2")

	pushNode := testOnlinePushNode(pathCfg+"_push2", pathDB+"_push2")
	if pushNode == nil {
		return nil
	}

	time.Sleep(200 * time.Millisecond)
	return pushNode
}

func testAllOnlineFree(node anonymity.INode, pathCfg, pathDB string) {
	defer func() {
		os.RemoveAll(pathCfg + "_push2")
		os.RemoveAll(pathDB + "_push2")
	}()
	types.StopAll([]types.ICommand{
		node,
		node.GetNetworkNode(),
	})
	types.CloseAll([]types.ICloser{
		node.GetWrapperDB(),
	})
}

func testOnlinePushNode(cfgPath, dbPath string) anonymity.INode {
	node := testRunNewNode(dbPath, testutils.TgAddrs[13])

	cfg, err := config.BuildConfig(cfgPath, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FKeySizeBits:      testutils.TcKeySize,
			FQueuePeriodMS:    testutils.TCQueuePeriod,
		},
	})
	if err != nil {
		return nil
	}

	node.HandleFunc(
		pkg_settings.CServiceMask,
		HandleServiceTCP(
			cfg,
			logger.NewLogger(
				logger.NewSettings(&logger.SSettings{}),
				func(_ logger.ILogArg) string {
					return ""
				},
			),
		),
	)
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024).GetPubKey())

	if err := node.GetNetworkNode().Run(); err != nil {
		return nil
	}
	return node
}
