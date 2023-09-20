package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/filesystem"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDB              = "database_test.db"
	tcPathConfig          = "config_test.cfg"
)

var (
	tcConfig = fmt.Sprintf(
		`{
			"settings": {
				"message_size_bytes": 8192,
				"work_size_bits": 20,
				"key_size_bits": %d,
				"queue_period_ms": 1000,
				"limit_void_size_bytes": 4096
			},
			"address": {
				"tcp": "test_address_tcp",
				"http": "test_address_http"
			},
			"network_key": "test",
			"connections": [
				"test_connect1",
				"test_connect2",
				"test_connect3"
			],
			"friends": {
				"test_recvr": "%s",
				"test_name1": "%s",
				"test_name2": "%s"
			},
			"services": {
				"test_service1": "test_address1",
				"test_service2": "test_address2",
				"test_service3": "test_address3"
			}
		}`,
		testutils.TcKeySize,
		testutils.TgPubKeys[0],
		testutils.TgPubKeys[1],
		testutils.TgPubKeys[2],
	)
)

func testStartServerHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:    addr,
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
		FMessage string `json:"message"`
	}

	var resp struct {
		FEcho  string `json:"echo"`
		FError int    `json:"error"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.FError = 1
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.FEcho = req.FMessage
	json.NewEncoder(w).Encode(resp)
}

func testAllCreate(cfgPath, dbPath, srvAddr string) (config.IWrapper, anonymity.INode, *http.Server) {
	wcfg := testNewWrapper(cfgPath)
	node := testRunNewNode(dbPath, "")
	srvc := testRunService(wcfg, node, srvAddr)
	time.Sleep(200 * time.Millisecond)
	return wcfg, node, srvc
}

func testAllFree(node anonymity.INode, srv *http.Server) {
	defer func() {
		os.RemoveAll(tcPathDB)
		os.RemoveAll(tcPathConfig)
	}()
	types.StopAll([]types.ICommand{
		node,
		node.GetNetworkNode(),
	})
	types.CloseAll([]types.ICloser{
		srv,
		node.GetWrapperDB(),
	})
}

func testRunService(wcfg config.IWrapper, node anonymity.INode, addr string) *http.Server {
	mux := http.NewServeMux()

	keySize := wcfg.GetConfig().GetSettings().GetKeySizeBits()
	ephPrivKey := asymmetric.NewRSAPrivKey(keySize)

	logger := logger.NewLogger(logger.NewSettings(&logger.SSettings{}))

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleConfigSettingsPath, HandleConfigSettingsAPI(wcfg, logger))
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, HandleConfigConnectsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, HandleConfigFriendsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, HandleNetworkOnlineAPI(logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, HandleNetworkRequestAPI(wcfg, logger, node, ephPrivKey))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, HandleNetworkMessageAPI(logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkKeyPath, HandleNetworkKeyAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNodeKeyPath, HandleNodeKeyAPI(wcfg, logger, node, ephPrivKey))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testNewWrapper(cfgPath string) config.IWrapper {
	filesystem.OpenFile(cfgPath).Write([]byte(tcConfig))
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	return config.NewWrapper(cfg)
}

func testRunNewNode(dbPath, addr string) anonymity.INode {
	node := testNewNode(dbPath, addr).HandleFunc(pkg_settings.CServiceMask, nil)
	if err := node.Run(); err != nil {
		return nil
	}
	return node
}

func testNewNode(dbPath, addr string) anonymity.INode {
	db, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath:     dbPath,
			FWorkSize: testutils.TCWorkSize,
			FPassword: "CIPHER",
		}),
	)
	if err != nil {
		return nil
	}
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:   "TEST",
			FRetryEnqueue:  0,
			FNetworkMask:   1,
			FFetchTimeWait: time.Minute,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		anonymity.NewWrapperDB().Set(db),
		testNewNetworkNode(addr),
		queue.NewMessageQueue(
			queue.NewSettings(&queue.SSettings{
				FMainCapacity: testutils.TCQueueCapacity,
				FPoolCapacity: testutils.TCQueueCapacity,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FWorkSizeBits:     testutils.TCWorkSize,
					FMessageSizeBytes: testutils.TCMessageSize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
			),
		),
		asymmetric.NewListPubKeys(),
	)
	return node
}

func testNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FCapacity:     testutils.TCCapacity,
			FMaxConnects:  testutils.TCMaxConnects,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		}),
	)
}
