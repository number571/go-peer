package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/app/state"
	"github.com/number571/go-peer/cmd/hlm/web"
)

func NotFoundPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"page404.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, s.GetTemplate())
	}
}
