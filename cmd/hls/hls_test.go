package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/client"
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
	tcPathDB     = "hls_test.db"
	tcPathConfig = "hls_test.cfg"
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDBTemplate      = "database_test_%d.db"
)

// client -> HLS -> server --\
// client <- HLS <- server <-/
func TestHLS(t *testing.T) {
	defer func() {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, 0))
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, 1))
		os.Remove(tcPathConfig)
	}()

	// server
	srv := testStartServerHTTP(t)
	defer srv.Close()

	// service
	node, err := testStartNodeHLS(t)
	if err != nil {
		t.Error(err)
		return
	}
	defer node.Close()

	// client
	nodeClient, err := testStartClientHLS()
	if err != nil {
		t.Error(err)
		return
	}
	defer nodeClient.Close()
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

func testStartNodeHLS(t *testing.T) (anonymity.INode, error) {
	node := testNewNode(0)
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}

	node.Handle(hls_settings.CHeaderHLS, testRouteHLS)
	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		err := node.Network().Listen(testutils.TgAddrs[4])
		if err != nil {
			t.Error(err)
		}
	}()

	return node, nil
}

func testRouteHLS(node anonymity.INode, sender asymmetric.IPubKey, pld payload.IPayload) []byte {
	// for test
	mapping := map[string]string{
		tcServiceAddressInHLS: testutils.TgAddrs[5],
	}

	// load request from message's body
	requestBytes := pld.Body()
	request := hls_network.LoadRequest(requestBytes)
	if request == nil {
		return nil
	}

	// get service's address by name
	address, ok := mapping[request.Host()]
	if !ok {
		return nil
	}

	// generate new request to serivce
	req, err := http.NewRequest(
		request.Method(),
		fmt.Sprintf("http://%s%s", address, request.Path()),
		bytes.NewReader(request.Body()),
	)
	if err != nil {
		return nil
	}

	req.Header.Add(hls_settings.CHeaderPubKey, sender.String())
	for key, val := range request.Head() {
		req.Header.Add(key, val)
	}

	// send request to service
	// and receive response from service
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// send result to client
	return data
}

// CLIENT

func testStartClientHLS() (anonymity.INode, error) {
	time.Sleep(200 * time.Millisecond)

	node := testNewNode(1)
	if node == nil {
		return nil, fmt.Errorf("node is not running")
	}

	node.Handle(hls_settings.CHeaderHLS, nil)
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

func testNewNode(i int) anonymity.INode {
	msgSize := uint64(100 << 10)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FTimeWait: 30 * time.Second,
		}),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:      fmt.Sprintf(tcPathDBTemplate, i),
				FHashing:   true,
				FCipherKey: []byte(testutils.TcKey1),
			}),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FCapacity:    (1 << 10),
				FMessageSize: msgSize,
				FMaxConnects: 10,
				FTimeWait:    5 * time.Second,
			}),
		),
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
	)
	if err := node.Run(); err != nil {
		return nil
	}
	return node
}
