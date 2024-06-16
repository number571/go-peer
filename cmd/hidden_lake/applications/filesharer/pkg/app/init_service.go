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
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/logger"
)

func (p *sApp) initIncomingServiceHTTP(pCtx context.Context, pHlsClient hls_client.IClient) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlf_settings.CLoadPath,
		handler.HandleIncomigLoadHTTP(pCtx, p.fHTTPLogger, p.fStgPath, pHlsClient),
	) // POST

	mux.HandleFunc(
		hlf_settings.CListPath,
		handler.HandleIncomigListHTTP(p.fHTTPLogger, p.fConfig, p.fStgPath),
	) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}

func (p *sApp) initInterfaceServiceHTTP(pCtx context.Context, pHlsClient hls_client.IClient) {
	mux := http.NewServeMux()
	mux.Handle(hlf_settings.CStaticPath, http.StripPrefix(
		hlf_settings.CStaticPath,
		handleFileServer(p.fHTTPLogger, p.fConfig, http.FS(web.GetStaticPath()))),
	)

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hlf_settings.CHandleIndexPath, handler.IndexPage(p.fHTTPLogger, p.fConfig))                              // GET, POST
	mux.HandleFunc(hlf_settings.CHandleAboutPath, handler.AboutPage(p.fHTTPLogger, p.fConfig))                              // GET
	mux.HandleFunc(hlf_settings.CHandleSettingsPath, handler.SettingsPage(pCtx, p.fHTTPLogger, cfgWrapper, pHlsClient))     // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hlf_settings.CHandleFriendsPath, handler.FriendsPage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient))        // GET, POST, DELETE
	mux.HandleFunc(hlf_settings.CHandleFriendsStoragePath, handler.StoragePage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient)) // GET, POST, DELETE

	p.fIntServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetInterface(),
		Handler:     mux, // http.TimeoutHandler returns bug with progress bar of file downloading
		ReadTimeout: (5 * time.Second),
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
