package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleStoragePointerPath, handler.HandlePointerAPI(p.fStorage, p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleStorageHashesPath, handler.HandleHashesAPI(p.fStorage, p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, handler.HandleMessageAPI(pCtx, p.fConfig, p.fStorage, p.fHTTPLogger, p.fAnonLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleConfigSettings, handler.HandleConfigSettingsAPI(p.fConfig, p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
