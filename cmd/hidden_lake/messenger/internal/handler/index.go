package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
)

func IndexPage(pStateManager state.IStateManager) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/" {
			NotFoundPage(pStateManager)(pW, pR)
			return
		}
		http.Redirect(pW, pR, "/about", http.StatusFound)
	}
}
