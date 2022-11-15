package handler

import "net/http"

func FaviconPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/favicon.ico" {
		NotFoundPage(w, r)
		return
	}
	http.Redirect(w, r, "/static/img/favicon.ico", http.StatusFound)
}
