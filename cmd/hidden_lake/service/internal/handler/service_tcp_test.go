package handler

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/settings"
	testutils "github.com/number571/go-peer/test/_data"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	"github.com/number571/go-peer/pkg/types"
)

const (
	tcPathDBTemplate = "database_test_tcp_%d.db"
)

func testCleanHLS() {
	os.RemoveAll(tcPathConfig + "_tcp")
	for i := 0; i < 2; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i))
	}
}

// client -> HLS -> server --\
// client <- HLS <- server <-/
func TestHLS(t *testing.T) {
	testCleanHLS()
	defer testCleanHLS()

	// server
	srv := testStartServerHTTP(testutils.TgAddrs[5])
	defer srv.Close()

	// service
	nodeService, err := testStartNodeHLS(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		types.StopAll([]types.ICommand{
			nodeService,
			nodeService.GetNetworkNode(),
		})
		types.CloseAll([]types.ICloser{
			nodeService.GetWrapperDB(),
		})
	}()

	// client
	nodeClient, err := testStartClientHLS()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		types.StopAll([]types.ICommand{
			nodeClient,
			nodeClient.GetNetworkNode(),
		})
		types.CloseAll([]types.ICloser{
			nodeClient.GetWrapperDB(),
		})
	}()
}

// HLS

func testStartNodeHLS(t *testing.T) (anonymity.INode, error) {
	rawCFG := &config.SConfig{
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[5],
		},
		SConfigSettings: settings.SConfigSettings{
			FSettings: settings.SConfigSettingsBlock{
				FMessageSize: testutils.TCMessageSize,
				FWorkSize:    testutils.TCWorkSize,
			},
		},
	}

	cfg, err := config.BuildConfig(tcPathConfig+"_tcp", rawCFG)
	if err != nil {
		return nil, err
	}

	node := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 0), testutils.TgAddrs[4])
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}

	node.HandleFunc(
		pkg_settings.CServiceMask,
		HandleServiceTCP(
			cfg,
			logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		),
	)
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).GetPubKey())

	if err := node.GetNetworkNode().Run(); err != nil {
		t.Error(err)
		return nil, nil
	}
	return node, nil
}

// CLIENT

func testStartClientHLS() (anonymity.INode, error) {
	time.Sleep(time.Second)

	node := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 1), "")
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}
	node.GetListPubKeys().AddPubKey(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).GetPubKey())

	if err := node.GetNetworkNode().AddConnect(testutils.TgAddrs[4]); err != nil {
		return nil, err
	}

	pld := adapters.NewPayload(
		pkg_settings.CServiceMask,
		request.NewRequest(http.MethodGet, tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
			ToBytes(),
	)

	pubKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).GetPubKey()
	respBytes, err := node.FetchPayload(pubKey, pld)
	if err != nil {
		return node, err
	}

	resp, err := response.LoadResponse(respBytes)
	if err != nil {
		return node, err
	}

	body := resp.GetBody()
	if string(body) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		return node, fmt.Errorf("result does not match; got '%s'", string(body))
	}

	return node, nil
}
