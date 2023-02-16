package app

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/handler"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
	"golang.org/x/net/websocket"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

const (
	initStart = time.Second * 3
)

var (
	_ types.IApp = &sApp{}
)

type sApp struct {
	fState          state.IState
	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	stg storage.IKeyValueStorage,
	wDB database.IWrapperDB,
	hlsClient hls_client.IClient,
	hltClient hlt_client.IClient,
) types.IApp {
	state := state.NewState(stg, wDB, hlsClient, hltClient)
	return &sApp{
		fState:          state,
		fIntServiceHTTP: initInterfaceServiceHTTP(cfg, state),
		fIncServiceHTTP: initIncomingServiceHTTP(cfg, state),
	}
}

func (app *sApp) Run() error {
	res := make(chan error)

	go func() {
		err := app.fIntServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	go func() {
		err := app.fIncServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			res <- err
			return
		}
	}()

	select {
	case err := <-res:
		app.Close()
		return err
	case <-time.After(initStart):
		return nil
	}
}

func (app *sApp) Close() error {
	lastErr := closer.CloseAll([]types.ICloser{
		app.fIntServiceHTTP,
		app.fIncServiceHTTP,
	})

	db := app.fState.GetWrapperDB().Get()
	if db != nil {
		if err := db.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func initIncomingServiceHTTP(
	cfg config.IConfig,
	state state.IState,
) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(hlm_settings.CPushPath, handler.HandleIncomigHTTP(state)) // POST

	return &http.Server{
		Addr:    cfg.Address().Incoming(),
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
		Addr:    cfg.Address().Interface(),
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
