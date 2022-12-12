package handler

import (
	"html/template"
	"net/http"

	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
)

func AboutPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		NotFoundPage(w, r)
		return
	}

	t, err := template.ParseFiles(
		hlm_settings.CPathTemplates+"index.html",
		hlm_settings.CPathTemplates+"about.html",
	)
	if err != nil {
		panic("can't load hmtl files")
	}
	t.Execute(w, nil)
}
