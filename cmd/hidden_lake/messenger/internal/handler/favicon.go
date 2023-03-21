package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
)

func FaviconPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/favicon.ico" {
			NotFoundPage(pState)(pW, pR)
			return
		}
		http.Redirect(pW, pR, "/static/img/favicon.ico", http.StatusFound)
	}
}
