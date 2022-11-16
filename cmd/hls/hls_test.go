package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/handler"
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/network/anonymity"
	anon_testutils "github.com/number571/go-peer/modules/network/anonymity/testutils"
	"github.com/number571/go-peer/modules/payload"
	"github.com/number571/go-peer/settings/testutils"
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDBTemplate      = "database_test_%d.db"
	tcPathConfig          = "config_test.cfg"
)

// client -> HLS -> server --\
// client <- HLS <- server <-/
func TestHLS(t *testing.T) {
	defer func() {
		os.RemoveAll(tcPathConfig)
		for i := 0; i < 2; i++ {
			os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, i))
		}
	}()

	// server
	srv := testStartServerHTTP(t)
	defer srv.Close()

	// service
	nodeService, err := testStartNodeHLS(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer closer.CloseAll([]modules.ICloser{
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
	defer closer.CloseAll([]modules.ICloser{
		nodeClient.KeyValueDB(),
		nodeClient.Network(),
		nodeClient,
	})
}

// SERVER

func testStartServerHTTP(t *testing.T) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testutils.TestEchoPage)

	srv := &http.Server{
		Addr:    testutils.TgAddrs[5],
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

// HLS

func testStartNodeHLS(t *testing.T) (anonymity.INode, error) {
	rawCFG := &config.SConfig{
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[5],
		},
	}

	cfg, err := config.NewConfig(tcPathConfig, rawCFG)
	if err != nil {
		return nil, err
	}

	node := testRunNewNode(0)
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}

	node.Handle(hls_settings.CHeaderHLS, handler.HandleServiceTCP(cfg))
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

	node := testRunNewNode(1)
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	conn := node.Network().Connect(testutils.TgAddrs[4])
	if conn == nil {
		return node, fmt.Errorf("conn is nil")
	}

	msg := payload.NewPayload(
		uint64(hls_settings.CHeaderHLS),
		hls_network.NewRequest("GET", tcServiceAddressInHLS, "/echo").
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
		return node, fmt.Errorf("result does not match; get '%s'", string(res))
	}

	return node, nil
}

func testRunNewNode(i int) anonymity.INode {
	node := anon_testutils.TestNewNode(fmt.Sprintf(tcPathDBTemplate, i))
	node.Handle(hls_settings.CHeaderHLS, nil)
	if err := node.Run(); err != nil {
		return nil
	}
	return node
}
