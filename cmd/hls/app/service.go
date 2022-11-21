package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/handler"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/network/anonymity"
)

func initServiceHTTP(wrapper config.IWrapper, node anonymity.INode) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(hls_settings.CHandleConfigConnects, handler.HandleConnectionsAPI(wrapper, node))
	mux.HandleFunc(hls_settings.CHandleConfigFriends, handler.HandleFriendsAPI(wrapper, node))
	mux.HandleFunc(hls_settings.CHandleNetworkOnline, handler.HandleOnlineAPI(node))
	mux.HandleFunc(hls_settings.CHandleNetworkPush, handler.HandlePushAPI(node))
	mux.HandleFunc(hls_settings.CHandleNodePubkey, handler.HandlePubKeyAPI(node))

	return &http.Server{
		Addr:    wrapper.Config().Address().HTTP(),
		Handler: mux,
	}
}
