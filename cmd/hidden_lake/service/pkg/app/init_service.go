package app

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fWrapper, p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkKeyPath, handler.HandleNetworkKeyAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNodeKeyPath, handler.HandleNodeKeyAPI(p.fWrapper, p.fHTTPLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fWrapper.GetConfig().GetAddress().GetHTTP(),
		Handler: mux,
	}
}
