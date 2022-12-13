package handler

import (
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hlm_settings "github.com/number571/go-peer/cmd/hlm/internal/settings"
)

func AboutPage(wDB database.IWrapperDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/about" {
			NotFoundPage(db)(w, r)
			return
		}
		t, err := template.ParseFiles(
			hlm_settings.CPathTemplates+"index.html",
			hlm_settings.CPathTemplates+"about.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}
		t.Execute(w, newTemplateData(db))
	}
}
