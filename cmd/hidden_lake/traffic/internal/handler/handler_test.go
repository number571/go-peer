package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/types"
)

const (
	databaseTemplate = "database_test_%s.db"
)

func testAllRun(addr string) (*http.Server, database.IKeyValueDB, hlt_client.IClient) {
	db := database.NewKeyValueDB(
		fmt.Sprintf(databaseTemplate, addr),
		database.NewSettings(&database.SSettings{
			FMessageSize: hlt_settings.CMessageSize,
			FWorkSize:    hlt_settings.CWorkSize,
		}),
	)

	srv := testRunService(db, addr)

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(fmt.Sprintf("http://%s", addr)),
	)
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

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, HandleHashesAPI(db))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, HandleMessageAPI(db))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}
