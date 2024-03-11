package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/handler"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/web"
	"github.com/number571/go-peer/pkg/logger"
)

func (p *sApp) initIncomingServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlf_settings.CLoadPath,
		handler.HandleIncomigLoadHTTP(pCtx, p.fHTTPLogger, p.fConfig, p.fStgPath),
	) // POST

	mux.HandleFunc(
		hlf_settings.CListPath,
		handler.HandleIncomigListHTTP(p.fHTTPLogger, p.fConfig, p.fStgPath),
	) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}

func (p *sApp) initInterfaceServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	mux.Handle(hlf_settings.CStaticPath, http.StripPrefix(
		hlf_settings.CStaticPath,
		handleFileServer(p.fHTTPLogger, p.fConfig, http.FS(web.GetStaticPath()))),
	)

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hlf_settings.CHandleIndexPath, handler.IndexPage(p.fHTTPLogger, p.fConfig))                  // GET, POST
	mux.HandleFunc(hlf_settings.CHandleAboutPath, handler.AboutPage(p.fHTTPLogger, p.fConfig))                  // GET
	mux.HandleFunc(hlf_settings.CHandleSettingsPath, handler.SettingsPage(pCtx, p.fHTTPLogger, cfgWrapper))     // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hlf_settings.CHandleFriendsPath, handler.FriendsPage(pCtx, p.fHTTPLogger, p.fConfig))        // GET, POST, DELETE
	mux.HandleFunc(hlf_settings.CHandleFriendsStoragePath, handler.StoragePage(pCtx, p.fHTTPLogger, p.fConfig)) // GET, POST, DELETE

	p.fIntServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetInterface(),
		ReadTimeout: (5 * time.Second),
		Handler:     mux, // http.TimeoutHandler send panic from websocket use
	}
}

func handleFileServer(pLogger logger.ILogger, pCfg config.IConfig, pFS http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := pFS.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(pLogger, pCfg)(w, r)
			return
		}
		http.FileServer(pFS).ServeHTTP(w, r)
	})
}
