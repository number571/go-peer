package app

import (
	"net/http"
	"os"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/handler"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"golang.org/x/net/websocket"
)

func (p *sApp) initIncomingServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(hlm_settings.CPushPath, handler.HandleIncomigHTTP(p.fState)) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetIncoming(),
		Handler: mux,
	}
}

func (p *sApp) initInterfaceServiceHTTP() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(p.fState, http.FS(web.GetStaticPath()))),
	)

	mux.HandleFunc("/", handler.IndexPage(p.fState))                    // GET
	mux.HandleFunc("/sign/out", handler.SignOutPage(p.fState))          // GET
	mux.HandleFunc("/sign/in", handler.SignInPage(p.fState))            // GET, POST
	mux.HandleFunc("/sign/up", handler.SignUpPage(p.fState))            // GET, POST
	mux.HandleFunc("/favicon.ico", handler.FaviconPage(p.fState))       // GET
	mux.HandleFunc("/about", handler.AboutPage(p.fState))               // GET
	mux.HandleFunc("/settings", handler.SettingsPage(p.fState))         // GET, POST, DELETE
	mux.HandleFunc("/qr/public_key", handler.QRPublicKeyPage(p.fState)) // GET
	mux.HandleFunc("/friends", handler.FriendsPage(p.fState))           // GET, POST, DELETE
	mux.HandleFunc("/friends/chat", handler.FriendsChatPage(p.fState))  // GET, POST

	mux.Handle("/friends/chat/ws", websocket.Handler(handler.FriendsChatWS))

	p.fIntServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetInterface(),
		Handler: mux,
	}
}

func handleFileServer(pState state.IState, pFS http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := pFS.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(pState)(w, r)
			return
		}
		http.FileServer(pFS).ServeHTTP(w, r)
	})
}
