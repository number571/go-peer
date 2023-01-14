package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/web"
)

func AboutPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/about" {
			NotFoundPage(s)(w, r)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"about.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		t.Execute(w, s.GetTemplate())
	}
}
