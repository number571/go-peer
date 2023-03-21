package app

import (
	"net/http"
	"os"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/handler"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"golang.org/x/net/websocket"
)

func initIncomingServiceHTTP(
	pCfg config.IConfig,
	pState state.IState,
) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(hlm_settings.CPushPath, handler.HandleIncomigHTTP(pState)) // POST

	return &http.Server{
		Addr:    pCfg.GetAddress().GetIncoming(),
		Handler: mux,
	}
}

func initInterfaceServiceHTTP(
	pCfg config.IConfig,
	pState state.IState,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(pState, http.FS(web.GetStaticPath()))),
	)

	mux.HandleFunc("/", handler.IndexPage(pState))                    // GET
	mux.HandleFunc("/sign/out", handler.SignOutPage(pState))          // GET
	mux.HandleFunc("/sign/in", handler.SignInPage(pState))            // GET, POST
	mux.HandleFunc("/sign/up", handler.SignUpPage(pState))            // GET, POST
	mux.HandleFunc("/favicon.ico", handler.FaviconPage(pState))       // GET
	mux.HandleFunc("/about", handler.AboutPage(pState))               // GET
	mux.HandleFunc("/settings", handler.SettingsPage(pState))         // GET, POST, DELETE
	mux.HandleFunc("/qr/public_key", handler.QRPublicKeyPage(pState)) // GET
	mux.HandleFunc("/friends", handler.FriendsPage(pState))           // GET, POST, DELETE
	mux.HandleFunc("/friends/chat", handler.FriendsChatPage(pState))  // GET, POST

	mux.Handle("/friends/chat/ws", websocket.Handler(handler.FriendsChatWS))

	return &http.Server{
		Addr:    pCfg.GetAddress().GetInterface(),
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
