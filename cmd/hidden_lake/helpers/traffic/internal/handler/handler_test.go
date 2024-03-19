package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/database"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/closer"
	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	databaseTemplate = "database_test_%s.db"
)

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   testutils.TCNetworkKey,
		FWorkSizeBits: testutils.TCWorkSize,
	})
}

func testAllRun(addr string) (*http.Server, context.CancelFunc, database.IDBWrapper, hlt_client.IClient) {
	db, err := database.NewDatabase(
		database.NewSettings(&database.SSettings{
			FPath:             fmt.Sprintf(databaseTemplate, addr),
			FNetworkKey:       testutils.TCNetworkKey,
			FWorkSizeBits:     testutils.TCWorkSize,
			FMessagesCapacity: testutils.TCCapacity,
		}),
	)
	if err != nil {
		return nil, nil, nil, nil
	}

	wDB := database.NewDBWrapper().Set(db)
	srv, _, cancel := testRunService(wDB, addr, "")

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", addr),
			&http.Client{Timeout: time.Minute},
			testNetworkMessageSettings(),
		),
	)

	time.Sleep(200 * time.Millisecond)
	return srv, cancel, wDB, hltClient
}

func testAllFree(addr string, srv *http.Server, cancel context.CancelFunc, wDB database.IDBWrapper) {
	defer func() {
		os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))
	}()
	cancel()
	_ = closer.CloseAll([]types.ICloser{srv, wDB})
}

func testRunService(wDB database.IDBWrapper, addr string, addrNode string) (*http.Server, connkeeper.IConnKeeper, context.CancelFunc) {
	mux := http.NewServeMux()

	connKeeperSettings := &connkeeper.SSettings{
		FDuration: time.Minute,
		FConnections: func() []string {
			return nil
		},
	}

	if addrNode != "" {
		connKeeperSettings.FConnections = func() []string {
			return []string{addrNode}
		}
	}

	connKeeper := connkeeper.NewConnKeeper(
		connkeeper.NewSettings(connKeeperSettings),
		testNewNetworkNode("").HandleFunc(
			1, // default value
			func(_ context.Context, _ network.INode, _ conn.IConn, _ net_message.IMessage) error {
				// pass response actions
				return nil
			},
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = connKeeper.Run(ctx) }()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes:   testutils.TCMessageSize,
			FWorkSizeBits:       testutils.TCWorkSize,
			FQueuePeriodMS:      hls_settings.CDefaultQueuePeriod,
			FLimitVoidSizeBytes: hls_settings.CDefaultLimitVoidSize,
			FKeySizeBits:        testutils.TcKeySize,
			FNetworkKey:         testutils.TCNetworkKey,
			FMessagesCapacity:   testutils.TCCapacity,
		},
	}

	node := network.NewNode(
		network.NewSettings(&network.SSettings{
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FNetworkKey:            testutils.TCNetworkKey,
				FWorkSizeBits:          testutils.TCWorkSize,
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: testutils.TCCapacity,
			}),
		),
	)

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleStoragePointerPath, HandlePointerAPI(wDB, logger))
	mux.HandleFunc(pkg_settings.CHandleStorageHashesPath, HandleHashesAPI(wDB, logger))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, HandleMessageAPI(ctx, cfg, wDB, logger, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigSettings, HandleConfigSettingsAPI(cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: time.Second,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}

	go func() { _ = srv.ListenAndServe() }()

	return srv, connKeeper, cancel
}

func testNewClient() client.IClient {
	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		privKey,
	)
}

func testNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FNetworkKey:            testutils.TCNetworkKey,
				FWorkSizeBits:          testutils.TCWorkSize,
				FLimitMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: testutils.TCCapacity,
			}),
		),
	)
}
