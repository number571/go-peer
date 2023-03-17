package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
)

func initServiceHTTP(cfg config.IConfig, connKeeper conn_keeper.IConnKeeper, wDB database.IWrapperDB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, handler.HandleHashesAPI(wDB))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, handler.HandleMessageAPI(connKeeper, wDB))

	return &http.Server{
		Addr:    cfg.GetAddress(),
		Handler: mux,
	}
}
