package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func initServiceHTTP(cfg config.IConfig, db database.IKeyValueDB) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleHashesPath, handler.HandleHashesAPI(db))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, handler.HandleMessageAPI(db))

	return &http.Server{
		Addr:    cfg.Address(),
		Handler: mux,
	}
}
