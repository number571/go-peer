package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
)

func IndexPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/" {
			NotFoundPage(pState)(pW, pR)
			return
		}
		http.Redirect(pW, pR, "/about", http.StatusFound)
	}
}
