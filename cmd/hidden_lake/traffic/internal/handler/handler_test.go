package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	databaseTemplate = "database_test_%s.db"
)

func testAllRun(addr, addrNode string) (*http.Server, conn_keeper.IConnKeeper, database.IWrapperDB, hlt_client.IClient) {
	db, err := database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{
			FPath:             fmt.Sprintf(databaseTemplate, addr),
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FCapacity:         testutils.TCCapacity,
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
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FWorkSizeBits:     testutils.TCWorkSize,
			}),
		),
	)

	time.Sleep(200 * time.Millisecond)
	return srv, connKeeper, wDB, hltClient
}

func testAllFree(addr string, srv *http.Server, connKeeper conn_keeper.IConnKeeper, wDB database.IWrapperDB) {
	defer func() {
		os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))
	}()
	types.StopAll([]types.ICommand{connKeeper})
	types.CloseAll([]types.ICloser{srv, wDB})
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

	sett := anonymity.NewSettings(&anonymity.SSettings{
		FServiceName:   "TEST",
		FRetryEnqueue:  0,
		FNetworkMask:   1,
		FFetchTimeWait: time.Minute,
	})
	connKeeper := conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(connKeeperSettings),
		testNewNetworkNode("").HandleFunc(
			sett.GetNetworkMask(), // default value
			func(_ network.INode, _ conn.IConn, _ []byte) {
				// pass response actions
			},
		),
	)
	if err := connKeeper.Run(); err != nil {
		return nil, nil
	}

	logger := logger.NewLogger(logger.NewSettings(&logger.SSettings{}))

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleHashesPath, HandleHashesAPI(wDB, logger))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, HandleMessageAPI(wDB, logger))

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
	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	return client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
		}),
		privKey,
	)
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
				FFetchTimeWait:    1, // not used
			}),
		}),
	)
}
