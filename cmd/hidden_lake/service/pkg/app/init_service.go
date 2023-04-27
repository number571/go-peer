package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func initServiceHTTP(pWrapper config.IWrapper, pNode anonymity.INode) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pWrapper, pNode))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(pWrapper, pNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(pNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pNode))
	mux.HandleFunc(pkg_settings.CHandleNodeKeyPath, handler.HandleNodeKeyAPI(pNode))

	return &http.Server{
		Addr:    pWrapper.GetConfig().GetAddress().GetHTTP(),
		Handler: mux,
	}
}
