package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	keySize := p.fWrapper.GetConfig().GetSettings().GetKeySizeBits()
	ephPrivKey := asymmetric.NewRSAPrivKey(keySize)

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fWrapper, p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(p.fWrapper, p.fHTTPLogger, p.fNode, ephPrivKey))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, handler.HandleNetworkMessageAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNetworkKeyPath, handler.HandleNetworkKeyAPI(p.fWrapper, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(pkg_settings.CHandleNodeKeyPath, handler.HandleNodeKeyAPI(p.fWrapper, p.fHTTPLogger, p.fNode, ephPrivKey))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fWrapper.GetConfig().GetAddress().GetHTTP(),
		Handler: mux,
	}
}
