package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
)

func IndexPage(wDB database.IWrapperDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/" {
			NotFoundPage(db)(w, r)
			return
		}
		http.Redirect(w, r, "/about", http.StatusFound)
	}
}
