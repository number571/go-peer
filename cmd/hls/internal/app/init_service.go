package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hls/internal/config"
	"github.com/number571/go-peer/cmd/hls/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
	"github.com/number571/go-peer/pkg/network/anonymity"
)

func initServiceHTTP(wrapper config.IWrapper, node anonymity.INode) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndex, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleConfigConnects, handler.HandleConfigConnectsAPI(wrapper, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriends, handler.HandleConfigFriendsAPI(wrapper, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnline, handler.HandleNetworkOnlineAPI(node))
	mux.HandleFunc(pkg_settings.CHandleNetworkPush, handler.HandleNetworkPushAPI(node))
	mux.HandleFunc(pkg_settings.CHandleNodeKey, handler.HandleNodeKeyAPI(node))

	return &http.Server{
		Addr:    wrapper.Config().Address().HTTP(),
		Handler: mux,
	}
}
