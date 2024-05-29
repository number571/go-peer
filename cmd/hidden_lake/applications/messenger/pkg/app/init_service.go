package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/handler"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/msgbroker"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/web"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/logger"
	"golang.org/x/net/websocket"
)

func (p *sApp) initIncomingServiceHTTP(
	pCtx context.Context,
	pHlsClient hls_client.IClient,
	pMsgBroker msgbroker.IMessageBroker,
) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlm_settings.CPushPath,
		handler.HandleIncomigHTTP(pCtx, p.fHTTPLogger, p.fDatabase, pMsgBroker, pHlsClient),
	) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}

func (p *sApp) initInterfaceServiceHTTP(
	pCtx context.Context,
	pHlsClient hls_client.IClient,
	pMsgBroker msgbroker.IMessageBroker,
) {
	mux := http.NewServeMux()
	mux.Handle(hlm_settings.CStaticPath, http.StripPrefix(
		hlm_settings.CStaticPath,
		handleFileServer(p.fHTTPLogger, p.fConfig, http.FS(web.GetStaticPath()))),
	)

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hlm_settings.CHandleIndexPath, handler.IndexPage(p.fHTTPLogger, p.fConfig))                                            // GET, POST
	mux.HandleFunc(hlm_settings.CHandleAboutPath, handler.AboutPage(p.fHTTPLogger, p.fConfig))                                            // GET
	mux.HandleFunc(hlm_settings.CHandleSettingsPath, handler.SettingsPage(pCtx, p.fHTTPLogger, cfgWrapper, pHlsClient))                   // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hlm_settings.CHandleFriendsPath, handler.FriendsPage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient))                      // GET, POST, DELETE
	mux.HandleFunc(hlm_settings.CHandleFriendsChatPath, handler.FriendsChatPage(pCtx, p.fHTTPLogger, p.fConfig, p.fDatabase, pHlsClient)) // GET, POST, PUT
	mux.HandleFunc(hlm_settings.CHandleFriendsUploadPath, handler.FriendsUploadPage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient))          // GET

	mux.Handle(hlm_settings.CHandleFriendsChatWSPath, websocket.Handler(handler.FriendsChatWS(pMsgBroker)))

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
