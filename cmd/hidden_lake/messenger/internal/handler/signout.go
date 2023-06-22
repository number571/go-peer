package handler

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
)

func SignOutPage(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/sign/out" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		if err := pStateManager.CloseState(); err != nil {
			fmt.Fprint(pW, "error: clean hls_client data")
			return
		}

		http.Redirect(pW, pR, "/about", http.StatusFound)
	}
}
