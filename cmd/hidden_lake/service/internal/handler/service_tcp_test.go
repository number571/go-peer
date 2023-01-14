package handler

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/internal/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	testutils "github.com/number571/go-peer/test/_data"

	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/payload"
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
	defer closer.CloseAll([]types.ICloser{
		nodeService.KeyValueDB(),
		nodeService.Network(),
		nodeService,
	})

	// client
	nodeClient, err := testStartClientHLS()
	if err != nil {
		t.Error(err)
		return
	}
	defer closer.CloseAll([]types.ICloser{
		nodeClient.KeyValueDB(),
		nodeClient.Network(),
		nodeClient,
	})
}

// HLS

func testStartNodeHLS(t *testing.T) (anonymity.INode, error) {
	rawCFG := &config.SConfig{
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[5],
		},
	}

	cfg, err := config.NewConfig(tcPathConfig+"_tcp", rawCFG)
	if err != nil {
		return nil, err
	}

	node := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 0))
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}

	node.Handle(hls_settings.CHeaderHLS, HandleServiceTCP(cfg))
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		err := node.Network().Listen(testutils.TgAddrs[4])
		if err != nil {
			t.Error(err)
		}
	}()

	return node, nil
}

// CLIENT

func testStartClientHLS() (anonymity.INode, error) {
	time.Sleep(time.Second)

	node := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 1))
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	_, err := node.Network().Connect(testutils.TgAddrs[4])
	if err != nil {
		return nil, err
	}

	msg := payload.NewPayload(
		uint64(hls_settings.CHeaderHLS),
		request.NewRequest("GET", tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
			Bytes(),
	)

	pubKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey()
	res, err := node.Request(pubKey, msg)
	if err != nil {
		return node, err
	}

	if string(res) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		return node, fmt.Errorf("result does not match; got '%s'", string(res))
	}

	return node, nil
}
