package handler

import "net/http"

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFoundPage(w, r)
		return
	}
	http.Redirect(w, r, "/about", http.StatusFound)
}
