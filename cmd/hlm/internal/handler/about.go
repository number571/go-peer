package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/number571/go-peer/cmd/hlm/internal/database"
	"github.com/number571/go-peer/cmd/hlm/web"
)

func AboutPage(wDB database.IWrapperDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := wDB.Get()
		if r.URL.Path != "/about" {
			NotFoundPage(db)(w, r)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"about.html",
		)
		if err != nil {
			fmt.Println(err)
			panic("can't load hmtl files")
		}

		t.Execute(w, newTemplateData(db))
	}
}
