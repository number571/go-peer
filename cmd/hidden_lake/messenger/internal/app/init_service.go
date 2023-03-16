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
	cfg config.IConfig,
	state state.IState,
) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(hlm_settings.CPushPath, handler.HandleIncomigHTTP(state)) // POST

	return &http.Server{
		Addr:    cfg.GetAddress().GetIncoming(),
		Handler: mux,
	}
}

func initInterfaceServiceHTTP(
	cfg config.IConfig,
	state state.IState,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(state, http.FS(web.GetStaticPath()))),
	)

	mux.HandleFunc("/", handler.IndexPage(state))                    // GET
	mux.HandleFunc("/sign/out", handler.SignOutPage(state))          // GET
	mux.HandleFunc("/sign/in", handler.SignInPage(state))            // GET, POST
	mux.HandleFunc("/sign/up", handler.SignUpPage(state))            // GET, POST
	mux.HandleFunc("/favicon.ico", handler.FaviconPage(state))       // GET
	mux.HandleFunc("/about", handler.AboutPage(state))               // GET
	mux.HandleFunc("/settings", handler.SettingsPage(state))         // GET, POST, DELETE
	mux.HandleFunc("/qr/public_key", handler.QRPublicKeyPage(state)) // GET
	mux.HandleFunc("/friends", handler.FriendsPage(state))           // GET, POST, DELETE
	mux.HandleFunc("/friends/chat", handler.FriendsChatPage(state))  // GET, POST

	mux.Handle("/friends/chat/ws", websocket.Handler(handler.FriendsChatWS))

	return &http.Server{
		Addr:    cfg.GetAddress().GetInterface(),
		Handler: mux,
	}
}

func handleFileServer(state state.IState, fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(state)(w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}
