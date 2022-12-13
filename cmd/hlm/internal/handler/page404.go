package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
)

func NotFoundPage(db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFiles(
			hlm_settings.CPathTemplates+"index.html",
			hlm_settings.CPathTemplates+"page404.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, newTemplateData(db))
	}
}
