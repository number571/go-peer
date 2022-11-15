package handler

import (
	"html/template"
	"net/http"

	hlm_settings "github.com/number571/go-peer/cmd/hlm/settings"
)

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	t, err := template.ParseFiles(
		hlm_settings.CPathTemplates+"index.html",
		hlm_settings.CPathTemplates+"page404.html",
	)
	if err != nil {
		panic("can't load hmtl files")
	}
	t.Execute(w, nil)
}
