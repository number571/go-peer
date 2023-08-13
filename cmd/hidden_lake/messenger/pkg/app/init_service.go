package app

import (
	"net/http"
	"os"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/handler"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	pkg_client "github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/logger"
	"golang.org/x/net/websocket"
)

func (p *sApp) initIncomingServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlm_settings.CPushPath,
		handler.HandleIncomigHTTP(p.fStateManager, p.fLogger),
	) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:    p.fWrapper.GetConfig().GetAddress().GetIncoming(),
		Handler: mux,
	}
}

func (p *sApp) initInterfaceServiceHTTP() {
	mux := http.NewServeMux()
	mux.Handle(hlm_settings.CStaticPath, http.StripPrefix(
		hlm_settings.CStaticPath,
		handleFileServer(p.fStateManager, p.fLogger, http.FS(web.GetStaticPath()))),
	)

	msgSize := p.fWrapper.GetConfig().GetMessageSizeBytes()
	keySize := p.fWrapper.GetConfig().GetKeySizeBits()
	msgLimitBytes := pkg_client.GetMessageLimit(msgSize, keySize)
	msgLimitBase64 := msgLimitBytes - (msgLimitBytes / 4) // https://ru.wikipedia.org/wiki/Base64

	mux.HandleFunc(hlm_settings.CHandleIndexPath, handler.IndexPage(p.fStateManager, p.fLogger))                                 // GET
	mux.HandleFunc(hlm_settings.CHandleSignOutPath, handler.SignOutPage(p.fStateManager, p.fLogger))                             // GET
	mux.HandleFunc(hlm_settings.CHandleSignInPath, handler.SignInPage(p.fStateManager, p.fLogger))                               // GET, POST
	mux.HandleFunc(hlm_settings.CHandleSignUpPath, handler.SignUpPage(p.fStateManager, p.fLogger))                               // GET, POST
	mux.HandleFunc(hlm_settings.CHandleFaviconPath, handler.FaviconPage(p.fStateManager, p.fLogger))                             // GET
	mux.HandleFunc(hlm_settings.CHandleAboutPath, handler.AboutPage(p.fStateManager, p.fLogger))                                 // GET
	mux.HandleFunc(hlm_settings.CHandleSettingsPath, handler.SettingsPage(p.fStateManager, p.fWrapper.GetEditor(), p.fLogger))   // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hlm_settings.CHandleQRPublicKeyKeyPath, handler.QRPublicKeyPage(p.fStateManager, p.fLogger))                  // GET
	mux.HandleFunc(hlm_settings.CHandleFriendsPath, handler.FriendsPage(p.fStateManager, p.fLogger))                             // GET, POST, DELETE
	mux.HandleFunc(hlm_settings.CHandleFriendsChatPath, handler.FriendsChatPage(p.fStateManager, p.fLogger, msgLimitBase64))     // GET, POST, PUT
	mux.HandleFunc(hlm_settings.CHandleFriendsUploadPath, handler.FriendsUploadPage(p.fStateManager, p.fLogger, msgLimitBase64)) // GET

	mux.Handle(hlm_settings.CHandleFriendsChatWSPath, websocket.Handler(handler.FriendsChatWS))

	p.fIntServiceHTTP = &http.Server{
		Addr:    p.fWrapper.GetConfig().GetAddress().GetInterface(),
		Handler: mux,
	}
}

func handleFileServer(pStateManager state.IStateManager, pLogger logger.ILogger, pFS http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := pFS.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(pStateManager, pLogger)(w, r)
			return
		}
		http.FileServer(pFS).ServeHTTP(w, r)
	})
}
