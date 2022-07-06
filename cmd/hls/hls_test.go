package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	hlsnet "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/netanon"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/testutils"
	"github.com/number571/go-peer/utils"
)

var (
	tgSettings = settings.NewSettings()
)

const (
	tcPathDB     = "hls_test.db"
	tcPathConfig = "hls_test.cfg"
)

const (
	tcServiceInHLS = "hidden-echo-service"
)

const (
	configBody = `{
	"clean_cron": "0 0 * * *",
	"address": {
		"hls": "localhost:9571",
		"http": ""
	},
	"f2f_mode": {
		"status": false,
		"pub_keys": []
	},
	"online_checker": {
		"status": false,
		"pub_keys": []
	},
	"connections": [],
	"services": {
		"hidden-echo-service": {
			"redirect": true,
			"address": "localhost:8081"
		}
	}
}`
)

func testHlsDefaultInit(dbPath, configPath string) {
	os.RemoveAll(tcPathDB)
	utils.OpenFile(configPath).Write([]byte(configBody))

	gDB = database.NewKeyValueDB(dbPath)
	gConfig = config.NewConfig(configPath)
}

// client -> HLS -> server -\
// client <- HLS <- server -/
func TestHLS(t *testing.T) {
	testHlsDefaultInit(tcPathDB, tcPathConfig)
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

	info, ok := gConfig.GetService(tcServiceInHLS)
	if !ok {
		return nil
	}

	srv := &http.Server{
		Addr:    info.Address(),
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
	client := client.NewClient(tgSettings, privKey)

	node := netanon.NewNode(client).
		Handle(hls_settings.CHeaderHLS, routeHLS)

	node.F2F().Switch(gConfig.F2F().Status())
	for _, pubKey := range gConfig.F2F().PubKeys() {
		node.F2F().Append(pubKey)
	}

	for _, addr := range gConfig.Connections() {
		conn := node.Network().Connect(addr)
		if conn == nil {
			t.Error("conn is nil")
		}
	}

	go func() {
		err := node.Network().Listen(gConfig.Address().HLS())
		if err != nil {
			t.Error(err)
		}
	}()

	return node
}

// CLIENT

func testStartClientHLS() error {
	priv := asymmetric.NewRSAPrivKey(testutils.TcAKeySize)
	client := client.NewClient(tgSettings, priv)

	node := netanon.NewNode(client).
		Handle(hls_settings.CHeaderHLS, nil)

	conn := node.Network().Connect(gConfig.Address().HLS())
	if conn == nil {
		return fmt.Errorf("conn is nil")
	}

	msg := payload.NewPayload(
		hls_settings.CHeaderHLS,
		hlsnet.NewRequest("GET", tcServiceInHLS, "/echo").
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
