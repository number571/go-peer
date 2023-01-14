package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
)

func IndexPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			NotFoundPage(s)(w, r)
			return
		}
		http.Redirect(w, r, "/about", http.StatusFound)
	}
}
