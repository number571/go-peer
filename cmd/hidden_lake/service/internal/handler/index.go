package handler

import (
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response(w, pkg_settings.CErrorNone, pkg_settings.CTitlePattern)
	}
}
