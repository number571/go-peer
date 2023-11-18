package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/interrupt"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/queue_set"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
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

func testAllRun(addr, addrNode string) (*http.Server, conn_keeper.IConnKeeper, database.IWrapperDB, hlt_client.IClient) {
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

	wDB := database.NewWrapperDB().Set(db)
	srv, connKeeper := testRunService(wDB, addr, addrNode)

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", addr),
			&http.Client{Timeout: time.Minute},
			testNetworkMessageSettings(),
		),
	)

	time.Sleep(200 * time.Millisecond)
	return srv, connKeeper, wDB, hltClient
}

func testAllFree(addr string, srv *http.Server, connKeeper conn_keeper.IConnKeeper, wDB database.IWrapperDB) {
	defer func() {
		os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))
	}()
	interrupt.StopAll([]types.ICommand{connKeeper})
	interrupt.CloseAll([]types.ICloser{srv, wDB})
}

func testRunService(wDB database.IWrapperDB, addr string, addrNode string) (*http.Server, conn_keeper.IConnKeeper) {
	mux := http.NewServeMux()

	connKeeperSettings := &conn_keeper.SSettings{
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

	connKeeper := conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(connKeeperSettings),
		testNewNetworkNode("").HandleFunc(
			1, // default value
			func(_ network.INode, _ conn.IConn, _ net_message.IMessage) error {
				// pass response actions
				return nil
			},
		),
	)
	if err := connKeeper.Run(); err != nil {
		return nil, nil
	}

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes:   testutils.TCMessageSize,
			FWorkSizeBits:       testutils.TCWorkSize,
			FQueuePeriodMS:      hls_settings.CDefaultQueuePeriod,
			FLimitVoidSizeBytes: hls_settings.CDefaultLimitVoidSize,
		},
	}

	node := network.NewNode(
		network.NewSettings(&network.SSettings{
			FMaxConnects:  testutils.TCMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FNetworkKey:       testutils.TCNetworkKey,
				FWorkSizeBits:     testutils.TCWorkSize,
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		}),
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: testutils.TCCapacity,
			}),
		),
	)

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleHashesPath, HandleHashesAPI(wDB, logger))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, HandleMessageAPI(cfg, wDB, logger, logger, node))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv, connKeeper
}

func testNewClient() client.IClient {
	privKey := asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
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
				FNetworkKey:       testutils.TCNetworkKey,
				FWorkSizeBits:     testutils.TCWorkSize,
				FMessageSizeBytes: testutils.TCMessageSize,
				FWaitReadDeadline: time.Hour,
				FReadDeadline:     time.Minute,
				FWriteDeadline:    time.Minute,
			}),
		}),
		queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: testutils.TCCapacity,
			}),
		),
	)
}
