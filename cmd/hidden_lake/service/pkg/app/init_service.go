package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(p.fWrapper, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fWrapper, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, handler.HandleNetworkMessageAPI(p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNodeKeyPath, handler.HandleNodeKeyAPI(p.fWrapper, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fWrapper.GetConfig().GetAddress().GetHTTP(),
		Handler: mux,
	}
}
