package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/internal/handler"
	hld_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc(hld_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hld_settings.CHandleConfigSettings, handler.HandleConfigSettingsAPI(p.fConfig, p.fHTTPLogger))
	mux.HandleFunc(hld_settings.CHandleNetworkDistributePath, handler.HandleNetworkDistributeAPI(pCtx, p.fConfig, p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
