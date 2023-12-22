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

	mux.HandleFunc(hls_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fWrapper, p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkKeyPath, handler.HandleNetworkKeyAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNodeKeyPath, handler.HandleNodeKeyAPI(p.fWrapper, p.fHTTPLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fWrapper.GetConfig().GetAddress().GetHTTP(),
		ReadTimeout: time.Second,
		// FetchTimeout = max of time waiting in the all handlers
		Handler: http.TimeoutHandler(mux, hls_settings.CFetchTimeout, "timeout"),
	}
}
