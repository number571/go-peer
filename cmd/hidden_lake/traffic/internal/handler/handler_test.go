package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
	"github.com/number571/go-peer/pkg/types"
)

const (
	databaseTemplate = "database_test_%s.db"
)

func testAllRun(addr string) (*http.Server, database.IKeyValueDB, hlt_client.IClient) {
	db := database.NewKeyValueDB(
		database.NewSettings(&database.SSettings{
			FMessageSize: hlt_settings.CMessageSize,
			FWorkSize:    hlt_settings.CWorkSize,
		}),
		fmt.Sprintf(databaseTemplate, addr),
	)

	srv := testRunService(db, addr)

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(fmt.Sprintf("http://%s", addr)),
	)

	time.Sleep(200 * time.Millisecond)
	return srv, db, hltClient
}

func testAllFree(addr string, srv *http.Server, db database.IKeyValueDB) {
	defer func() {
		os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))
	}()
	closer.CloseAll([]types.ICloser{srv, db})
}

func testRunService(db database.IKeyValueDB, addr string) *http.Server {
	mux := http.NewServeMux()

	// TODO: make node with connection
	connKeeper := conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{}),
		network.NewNode(network.NewSettings(&network.SSettings{})),
	)

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, HandleHashesAPI(db))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, HandleMessageAPI(db))
	mux.HandleFunc(pkg_settings.CHandleBroadcastPath, HandleBroadcastAPI(db, connKeeper))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}
