package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/client"
	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/database"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/netanon"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/testutils"
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
	node := testStartNodeHLS(t)
	defer node.Close()

	// client
	nodeClient, err := testStartClientHLS()
	if err != nil {
		t.Error(err)
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

func testStartNodeHLS(t *testing.T) netanon.INode {
	node := testNewNode(0).
		Handle(hls_settings.CHeaderHLS, testRouteHLS)

	node.F2F().Append(asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey())

	go func() {
		err := node.Network().Listen(testutils.TgAddrs[4])
		if err != nil {
			t.Error(err)
		}
	}()

	return node
}

func testRouteHLS(node netanon.INode, _ asymmetric.IPubKey, pld payload.IPayload) []byte {
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
	addr, ok := mapping[request.Host()]
	if !ok {
		return nil
	}

	// generate new request to serivce
	req, err := http.NewRequest(
		request.Method(),
		fmt.Sprintf("http://%s%s", addr, request.Path()),
		bytes.NewReader(request.Body()),
	)
	if err != nil {
		return nil
	}
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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	// send result to client
	return data
}

// CLIENT

func testStartClientHLS() (netanon.INode, error) {
	node := testNewNode(1).
		Handle(hls_settings.CHeaderHLS, nil)

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
			ToBytes(),
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

func testNewNode(i int) netanon.INode {
	msgSize := uint64(100 << 10)
	return netanon.NewNode(
		netanon.NewSettings(&netanon.SSettings{
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
				FRetryNum:    2,
				FCapacity:    (1 << 10),
				FMessageSize: msgSize,
				FMaxConns:    10,
				FMaxMessages: 20,
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
}
