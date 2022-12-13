package app

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hlm/internal/config"
	"github.com/number571/go-peer/cmd/hlm/internal/database"
	"github.com/number571/go-peer/cmd/hlm/internal/handler"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
	"golang.org/x/net/websocket"
)

var (
	_ types.IApp = &sApp{}
)

type sApp struct {
	fStorage        storage.IKeyValueStorage
	fDatabase       database.IWrapperDB
	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	client hls_client.IClient,
) types.IApp {
	var wDB = database.NewWrapperDB()
	stg, err := initCryptoStorage(cfg)
	if err != nil {
		panic(err)
	}
	return &sApp{
		fStorage:        stg,
		fDatabase:       wDB,
		fIntServiceHTTP: initIntServiceHTTP(cfg, wDB, client, stg),
		fIncServiceHTTP: initIncServiceHTTP(cfg, wDB, client),
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
	case <-time.After(time.Second * 3):
		return nil
	}
}

func (app *sApp) Close() error {
	return closer.CloseAll([]types.ICloser{
		app.fIntServiceHTTP,
		app.fIncServiceHTTP,
		app.fDatabase,
	})
}

func initCryptoStorage(cfg config.IConfig) (storage.IKeyValueStorage, error) {
	return storage.NewCryptoStorage(
		hlm_settings.CPathSTG,
		[]byte(cfg.StorageKey()),
		hlm_settings.CWorkForKeys,
	)
}

func initIncServiceHTTP(
	cfg config.IConfig,
	wDB database.IWrapperDB,
	client hls_client.IClient,
) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", handler.HandleIncomigHTTP(wDB, client))

	return &http.Server{
		Addr:    cfg.Address().Incoming(),
		Handler: mux,
	}
}

func initIntServiceHTTP(
	cfg config.IConfig,
	wDB database.IWrapperDB,
	client hls_client.IClient,
	stg storage.IKeyValueStorage,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(wDB, http.Dir(hlm_settings.CPathStatic))),
	)

	mux.HandleFunc("/", handler.IndexPage(wDB))                            // GET
	mux.HandleFunc("/sign/out", handler.SignOutPage(wDB, client))          // GET
	mux.HandleFunc("/sign/in", handler.SignInPage(wDB, client, stg))       // GET, POST
	mux.HandleFunc("/sign/up", handler.SignUpPage(wDB, stg))               // GET, POST
	mux.HandleFunc("/favicon.ico", handler.FaviconPage(wDB))               // GET
	mux.HandleFunc("/about", handler.AboutPage(wDB))                       // GET
	mux.HandleFunc("/settings", handler.SettingsPage(wDB, client))         // GET, POST, DELETE
	mux.HandleFunc("/qr/public_key", handler.QRPublicKeyPage(wDB, client)) // GET
	mux.HandleFunc("/friends", handler.FriendsPage(wDB, client))           // GET, POST, DELETE
	mux.HandleFunc("/friends/chat", handler.FriendsChatPage(wDB, client))  // GET, POST

	mux.Handle("/friends/chat/ws", websocket.Handler(handler.FriendsChatWS))

	return &http.Server{
		Addr:    cfg.Address().Interface(),
		Handler: mux,
	}
}

func handleFileServer(wDB database.IWrapperDB, fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(wDB.Get())(w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}
