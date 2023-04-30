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
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"
	testutils "github.com/number571/go-peer/test/_data"
	anon_testutils "github.com/number571/go-peer/test/_data/anonymity"
)

const (
	databaseTemplate = "database_test_%s.db"
)

func testAllRun(addr, addrNode string) (*http.Server, conn_keeper.IConnKeeper, database.IWrapperDB, hlt_client.IClient) {
	db, err := database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{
			FPath:        fmt.Sprintf(databaseTemplate, addr),
			FMessageSize: anon_testutils.TCMessageSize,
			FWorkSize:    anon_testutils.TCWorkSize,
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
				FMessageSize: anon_testutils.TCMessageSize,
				FWorkSize:    anon_testutils.TCWorkSize,
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

	connKeeperSettings := &conn_keeper.SSettings{}
	if addrNode != "" {
		connKeeperSettings.FConnections = func() []string {
			return []string{addrNode}
		}
	}

	sett := anonymity.NewSettings(&anonymity.SSettings{})
	connKeeper := conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(connKeeperSettings),
		anon_testutils.TestNewNetworkNode("").HandleFunc(
			sett.GetNetworkMask(), // default value
			func(_ network.INode, _ conn.IConn, _ []byte) {
				// pass response actions
			},
		),
	)
	if err := connKeeper.Run(); err != nil {
		return nil, nil
	}

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, HandleHashesAPI(wDB))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, HandleMessageAPI(connKeeper, wDB))

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
			FMessageSize: anon_testutils.TCMessageSize,
			FWorkSize:    anon_testutils.TCWorkSize,
		}),
		privKey,
	)
}
