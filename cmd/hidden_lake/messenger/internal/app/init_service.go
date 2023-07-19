package app

import (
	"net/http"
	"os"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/handler"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	pkg_client "github.com/number571/go-peer/pkg/client"
	"golang.org/x/net/websocket"
)

func (p *sApp) initIncomingServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(hlm_settings.CPushPath, handler.HandleIncomigHTTP(p.fStateManager)) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetIncoming(),
		Handler: mux,
	}
}

func (p *sApp) initInterfaceServiceHTTP() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(p.fStateManager, http.FS(web.GetStaticPath()))),
	)

	msgSize := p.fConfig.GetMessageSizeBytes()
	keySize := p.fConfig.GetKeySizeBits()
	msgLimit := pkg_client.GetMessageLimit(msgSize, keySize)

	mux.HandleFunc("/", handler.IndexPage(p.fStateManager))                                 // GET
	mux.HandleFunc("/sign/out", handler.SignOutPage(p.fStateManager))                       // GET
	mux.HandleFunc("/sign/in", handler.SignInPage(p.fStateManager))                         // GET, POST
	mux.HandleFunc("/sign/up", handler.SignUpPage(p.fStateManager))                         // GET, POST
	mux.HandleFunc("/favicon.ico", handler.FaviconPage(p.fStateManager))                    // GET
	mux.HandleFunc("/about", handler.AboutPage(p.fStateManager))                            // GET
	mux.HandleFunc("/settings", handler.SettingsPage(p.fStateManager))                      // GET, POST, DELETE
	mux.HandleFunc("/qr/public_key", handler.QRPublicKeyPage(p.fStateManager))              // GET
	mux.HandleFunc("/friends", handler.FriendsPage(p.fStateManager))                        // GET, POST, DELETE
	mux.HandleFunc("/friends/chat", handler.FriendsChatPage(p.fStateManager, msgLimit))     // GET, POST, PUT
	mux.HandleFunc("/friends/upload", handler.FriendsUploadPage(p.fStateManager, msgLimit)) // GET

	mux.Handle("/friends/chat/ws", websocket.Handler(handler.FriendsChatWS))

	p.fIntServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetInterface(),
		Handler: mux,
	}
}

func handleFileServer(pStateManager state.IStateManager, pFS http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := pFS.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(pStateManager)(w, r)
			return
		}
		http.FileServer(pFS).ServeHTTP(w, r)
	})
}
