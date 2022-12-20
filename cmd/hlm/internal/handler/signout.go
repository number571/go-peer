package handler

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/app/state"
)

func SignOutPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sign/out" {
			NotFoundPage(s)(w, r)
			return
		}

		if !s.IsActive() {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		if err := s.ClearActiveState(); err != nil {
			fmt.Fprint(w, "error: clean hls_client data")
			return
		}

		http.Redirect(w, r, "/about", http.StatusFound)
	}
}
