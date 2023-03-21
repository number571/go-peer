package handler

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
)

func SignOutPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/sign/out" {
			NotFoundPage(pState)(pW, pR)
			return
		}

		if !pState.IsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		if err := pState.ClearActiveState(); err != nil {
			fmt.Fprint(pW, "error: clean hls_client data")
			return
		}

		http.Redirect(pW, pR, "/about", http.StatusFound)
	}
}
