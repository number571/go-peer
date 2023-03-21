package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

func NotFoundPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, _ *http.Request) {
		pW.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"page404.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(pW, pState.GetTemplate())
	}
}
