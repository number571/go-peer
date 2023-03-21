package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/pkg/network/conn_keeper"
)

func initServiceHTTP(pCfg config.IConfig, pConnKeeper conn_keeper.IConnKeeper, pWrapperDB database.IWrapperDB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, handler.HandleHashesAPI(pWrapperDB))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, handler.HandleMessageAPI(pConnKeeper, pWrapperDB))

	return &http.Server{
		Addr:    pCfg.GetAddress(),
		Handler: mux,
	}
}
