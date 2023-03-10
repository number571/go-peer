package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
)

func initServiceHTTP(cfg config.IConfig, connKeeper conn_keeper.IConnKeeper, db database.IKeyValueDB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, handler.HandleHashesAPI(db))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, handler.HandleMessageAPI(connKeeper, db))

	return &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: mux,
	}
}
