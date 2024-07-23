package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/closer"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache/lru"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDBTemplate      = "database_test_%d.db"
	tcPathConfigTemplate  = "config_test_%d.yml"
)

var (
	tcConfig = fmt.Sprintf(`settings:
  message_size_bytes: 8192
  work_size_bits: 22
  key_size_bits: %d
  fetch_timeout_ms: 60000
  queue_period_ms: 1000
  rand_message_size_bytes: 4096
  network_key: test
address:
  tcp: test_address_tcp
  http: test_address_http
connections:
  - test_connect1
  - test_connect2
  - test_connect3
friends:
  test_recvr: %s
  test_name1: %s
  test_name2: %s
services:
  test_service1: 
    host: test_address1
  test_service2: 
    host: test_address2
  test_service3: 
    host: test_address3
`,
		testutils.TcKeySize,
		testutils.TgPubKeys[0],
		testutils.TgPubKeys[1],
		testutils.TgPubKeys[2],
	)
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SHandlerError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testStartServerHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()

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
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp.FEcho = req.FMessage
	_ = json.NewEncoder(w).Encode(resp)
}

func testAllCreate(cfgPath, dbPath, srvAddr string) (config.IWrapper, anonymity.INode, context.Context, context.CancelFunc, *http.Server) {
	wcfg := testNewWrapper(cfgPath)
	node, ctx, cancel := testRunNewNode(dbPath, "")
	srvc := testRunService(ctx, wcfg, node, srvAddr)
	time.Sleep(200 * time.Millisecond)
	return wcfg, node, ctx, cancel, srvc
}

func testAllFree(node anonymity.INode, cancel context.CancelFunc, srv *http.Server, pathCfg, pathDB string) {
	defer func() {
		os.RemoveAll(pathDB)
		os.RemoveAll(pathCfg)
	}()
	cancel()
	_ = closer.CloseAll([]types.ICloser{
		srv,
		node.GetKVDatabase(),
		node.GetNetworkNode(),
	})
}

func testRunService(ctx context.Context, wcfg config.IWrapper, node anonymity.INode, addr string) *http.Server {
	mux := http.NewServeMux()

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	cfg := wcfg.GetConfig()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleConfigSettingsPath, HandleConfigSettingsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, HandleConfigConnectsAPI(ctx, wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, HandleConfigFriendsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, HandleNetworkOnlineAPI(logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, HandleNetworkRequestAPI(ctx, cfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkPubKeyPath, HandleNetworkPubKeyAPI(logger, node))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}

func testNewWrapper(cfgPath string) config.IWrapper {
	_ = os.WriteFile(cfgPath, []byte(tcConfig), 0o600)
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	return config.NewWrapper(cfg)
}

func testRunNewNode(dbPath, addr string) (anonymity.INode, context.Context, context.CancelFunc) {
	os.RemoveAll(dbPath)
	node := testNewNode(dbPath, addr).HandleFunc(pkg_settings.CServiceMask, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = node.Run(ctx) }()
	return node, ctx, cancel
}

func testNewNode(dbPath, addr string) anonymity.INode {
	db, err := database.NewKVDatabase(dbPath)
	if err != nil {
		panic(err)
	}
	networkMask := uint32(1)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  "TEST",
			FNetworkMask:  networkMask,
			FFetchTimeout: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		db,
		testNewNetworkNode(addr),
		queue.NewMessageQueueProcessor(
			queue.NewSettings(&queue.SSettings{
				FNetworkMask:      networkMask,
				FWorkSizeBits:     testutils.TCWorkSize,
				FMainPoolCapacity: testutils.TCQueueCapacity,
				FRandPoolCapacity: testutils.TCQueueCapacity,
				FQueuePeriod:      500 * time.Millisecond,
			}),
			queue.NewVSettings(&queue.SVSettings{}),
			client.NewClient(
				message.NewSettings(&message.SSettings{
					FMessageSizeBytes: testutils.TCMessageSize,
					FKeySizeBits:      testutils.TcKeySize,
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
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FWorkSizeBits:          testutils.TCWorkSize,
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		conn.NewVSettings(&conn.SVSettings{}),
		lru.NewLRUCache(testutils.TCCapacity),
	)
}
