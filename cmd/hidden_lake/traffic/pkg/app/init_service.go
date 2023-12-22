package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/handler"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc(pkg_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandlePointerPath, handler.HandlePointerAPI(p.fWrapperDB, p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleHashesPath, handler.HandleHashesAPI(p.fWrapperDB, p.fHTTPLogger))
	mux.HandleFunc(pkg_settings.CHandleMessagePath, handler.HandleMessageAPI(pCtx, p.fConfig, p.fWrapperDB, p.fHTTPLogger, p.fAnonLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}
