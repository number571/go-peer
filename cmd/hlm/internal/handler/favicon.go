package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/app/state"
)

func FaviconPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/favicon.ico" {
			NotFoundPage(s)(w, r)
			return
		}
		http.Redirect(w, r, "/static/img/favicon.ico", http.StatusFound)
	}
}
