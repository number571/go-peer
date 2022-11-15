package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hls/config"
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/payload"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage/database"
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
	db, nnode, nodeService, err := testStartNodeHLS(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer closer.CloseAll([]modules.ICloser{db, nnode, nodeService})

	// client
	db, nnode, nodeClient, err := testStartClientHLS()
	if err != nil {
		t.Error(err)
		return
	}
	defer closer.CloseAll([]modules.ICloser{db, nnode, nodeClient})
}

// SERVER

func testStartServerHTTP(t *testing.T) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:    testutils.TgAddrs[5],
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testEchoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Message string `json:"message"`
	}

	var resp struct {
		Echo  string `json:"echo"`
		Error int    `json:"error"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.Error = 1
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Echo = req.Message
	json.NewEncoder(w).Encode(resp)
}

// HLS

func testStartNodeHLS(t *testing.T) (database.IKeyValueDB, network.INode, anonymity.INode, error) {
	rawCFG := &config.SConfig{
		FAddress: &config.SAddress{
			FTCP: testutils.TgAddrs[4],
		},
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[5],
		},
	}

	cfg, err := config.NewConfig(tcPathConfig, rawCFG)
	if err != nil {
		return nil, nil, nil, err
	}

	db, nnode, node := testRunNewNode(0)
	if node == nil {
		return nil, nil, nil, fmt.Errorf("node is not running")
	}

	node.Handle(hls_settings.CHeaderHLS, handleServiceTCP(cfg))
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		err := node.Network().Listen(testutils.TgAddrs[4])
		if err != nil {
			t.Error(err)
		}
	}()

	return db, nnode, node, nil
}

// CLIENT

func testStartClientHLS() (database.IKeyValueDB, network.INode, anonymity.INode, error) {
	time.Sleep(time.Second)

	db, nnode, node := testRunNewNode(1)
	if node == nil {
		return nil, nil, nil, fmt.Errorf("node is not running")
	}
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	conn := node.Network().Connect(testutils.TgAddrs[4])
	if conn == nil {
		return db, nnode, node, fmt.Errorf("conn is nil")
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
		return db, nnode, node, err
	}

	if string(res) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		return db, nnode, node, fmt.Errorf("result does not match; get '%s'", string(res))
	}

	return db, nnode, node, nil
}

func testRunNewNode(i int) (database.IKeyValueDB, network.INode, anonymity.INode) {
	msgSize := uint64(100 << 10)
	db := database.NewLevelDB(
		database.NewSettings(&database.SSettings{
			FPath:      fmt.Sprintf(tcPathDBTemplate, i),
			FHashing:   true,
			FCipherKey: []byte(testutils.TcKey1),
		}),
	)
	nnode := network.NewNode(
		network.NewSettings(&network.SSettings{
			FCapacity:    (1 << 10),
			FMessageSize: msgSize,
			FMaxConnects: 10,
			FTimeWait:    5 * time.Second,
		}),
	)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FTimeWait: 30 * time.Second,
		}),
		db,
		nnode,
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     10,
				FPullCapacity: 5,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    10,
					FMessageSize: msgSize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		friends.NewF2F(),
	).Handle(hls_settings.CHeaderHLS, nil)
	if err := node.Run(); err != nil {
		return nil, nil, nil
	}
	return db, nnode, node
}
