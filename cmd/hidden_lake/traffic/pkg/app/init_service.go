package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fLogger))
	mux.HandleFunc(pkg_settings.CHandleHashesPath, handler.HandleHashesAPI(p.fWrapperDB, p.fLogger))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, handler.HandleMessageAPI(p.fWrapperDB, p.fLogger))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetHTTP(),
		Handler: mux,
	}
}
