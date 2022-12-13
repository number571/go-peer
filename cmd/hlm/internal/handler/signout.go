package handler

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
)

func SignOutPage(wDB database.IWrapperDB, client hls_client.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/sign/out" {
			NotFoundPage(db)(w, r)
			return
		}

		if db == nil {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		if err := wDB.Close(); err != nil {
			fmt.Fprint(w, "error: close database")
			return
		}

		http.Redirect(w, r, "/about", http.StatusFound)
	}
}
