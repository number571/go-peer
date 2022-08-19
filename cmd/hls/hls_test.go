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
	hlsnet "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/crypto/asymmetric"
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
)

// client -> HLS -> server -\
// client <- HLS <- server -/
func TestHLS(t *testing.T) {
	defer func() {
		os.RemoveAll(tcPathDB)
		os.Remove(tcPathConfig)
	}()

	// server
	srv := testStartServerHTTP(t)
	defer srv.Close()

	// service
	node := testStartNodeHLS(t)
	defer node.Close()

	// client
	err := testStartClientHLS()
	if err != nil {
		t.Error(err)
	}
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
	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	client := client.NewClient(client.NewSettings(10, (1<<20)), privKey)

	node := testNewNode(client).
		Handle(hls_settings.CHeaderHLS, testRouteHLS)

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
	request := hlsnet.LoadRequest(requestBytes)
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

func testStartClientHLS() error {
	privKey := asymmetric.NewRSAPrivKey(testutils.TcAKeySize)
	client := client.NewClient(client.NewSettings(10, (1<<20)), privKey)

	node := testNewNode(client).
		Handle(hls_settings.CHeaderHLS, nil)

	conn := node.Network().Connect(testutils.TgAddrs[4])
	if conn == nil {
		return fmt.Errorf("conn is nil")
	}

	msg := payload.NewPayload(
		uint64(hls_settings.CHeaderHLS),
		hlsnet.NewRequest("GET", tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
			ToBytes(),
	)

	pubKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey).PubKey()
	res, err := node.Request(pubKey, msg)
	if err != nil {
		return err
	}

	if string(res) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		return fmt.Errorf("result does not match; get '%s'", string(res))
	}

	return nil
}

func testNewNode(client client.IClient) netanon.INode {
	msgSize := uint64(1 << 20)
	return netanon.NewNode(
		netanon.NewSettings(
			3,
			20*time.Second,
		),
		client,
		network.NewNode(network.NewSettings(
			msgSize,
			10,
			1024,
			10,
			20,
			5*time.Second,
		)),
		queue.NewQueue(
			queue.NewSettings(
				10,
				5,
				300*time.Millisecond,
			),
			client,
		),
		friends.NewF2F(),
	)
}
