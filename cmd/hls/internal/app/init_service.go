package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hls/internal/config"
	"github.com/number571/go-peer/cmd/hls/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func initServiceHTTP(wrapper config.IWrapper, node anonymity.INode) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndex, handler.HandleIndexAPI())
	mux.HandleFunc(pkg_settings.CHandleConfigConnects, handler.HandleConnectionsAPI(wrapper, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriends, handler.HandleFriendsAPI(wrapper, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnline, handler.HandleOnlineAPI(node))
	mux.HandleFunc(pkg_settings.CHandleNetworkPush, handler.HandlePushAPI(node))
	mux.HandleFunc(pkg_settings.CHandleNodePubkey, handler.HandlePubKeyAPI(node))

	return &http.Server{
		Addr:    wrapper.Config().Address().HTTP(),
		Handler: mux,
	}
}
