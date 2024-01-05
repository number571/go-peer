package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	cfg := p.fCfgW.GetConfig()

	mux.HandleFunc(hls_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fCfgW, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, p.fCfgW, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fCfgW, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, cfg, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkPubKeyPath, handler.HandleNetworkPubKeyAPI(p.fHTTPLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:        cfg.GetAddress().GetHTTP(),
		ReadTimeout: (5 * time.Second),
		// FetchTimeout = max of time waiting in the all handlers
		Handler: http.TimeoutHandler(mux, hls_settings.CFetchTimeout, "timeout"),
	}
}
