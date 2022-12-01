package app

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hlm/config"
	"github.com/number571/go-peer/cmd/hlm/database"
	"github.com/number571/go-peer/cmd/hlm/handler"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
	"github.com/number571/go-peer/cmd/hls/hlc"
	"github.com/number571/go-peer/modules"
	"github.com/number571/go-peer/modules/closer"
	"golang.org/x/net/websocket"
)

var (
	_ IApp = &sApp{}
)

type sApp struct {
	fDB             database.IKeyValueDB
	fWebServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
}

func NewApp(
	cfg config.IConfig,
	client hlc.IClient,
	db database.IKeyValueDB,
) IApp {
	return &sApp{
		fDB:             db,
		fWebServiceHTTP: initWebServiceHTTP(cfg, client, db),
		fIncServiceHTTP: initIncServiceHTTP(cfg, db),
	}
}

func (app *sApp) Run() error {
	res := make(chan error)

	go func() {
		err := app.fWebServiceHTTP.ListenAndServe()
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
	return closer.CloseAll([]modules.ICloser{
		app.fWebServiceHTTP,
		app.fIncServiceHTTP,
		app.fDB,
	})
}

func initIncServiceHTTP(cfg config.IConfig, db database.IKeyValueDB) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", handler.HandleIncomigHTTP(db))

	return &http.Server{
		Addr:    cfg.Address().Incoming(),
		Handler: mux,
	}
}

func initWebServiceHTTP(cfg config.IConfig, client hlc.IClient, db database.IKeyValueDB) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(http.Dir(hlm_settings.CPathStatic))),
	)

	mux.HandleFunc("/", handler.IndexPage)                               // GET
	mux.HandleFunc("/favicon.ico", handler.FaviconPage)                  // GET
	mux.HandleFunc("/about", handler.AboutPage)                          // GET
	mux.HandleFunc("/settings", handler.SettingsPage(client))            // GET, POST, DELETE
	mux.HandleFunc("/qr/public_key", handler.QRPublicKeyPage(client))    // GET
	mux.HandleFunc("/friends", handler.FriendsPage(client))              // GET, POST, DELETE
	mux.HandleFunc("/friends/chat", handler.FriendsChatPage(client, db)) // GET, POST

	mux.Handle("/friends/chat/ws", websocket.Handler(handler.FriendsChatWS))

	return &http.Server{
		Addr:    cfg.Address().Interface(),
		Handler: mux,
	}
}

func handleFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}
