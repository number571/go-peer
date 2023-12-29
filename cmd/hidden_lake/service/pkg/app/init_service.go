package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hls_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(cfgWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, cfgWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(cfgWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, p.fConfig, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkPubKeyPath, handler.HandleNetworkPubKeyAPI(p.fHTTPLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		ReadTimeout: (5 * time.Second),
		// FetchTimeout = max of time waiting in the all handlers
		Handler: http.TimeoutHandler(mux, hls_settings.CFetchTimeout, "timeout"),
	}
}
