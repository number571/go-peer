package handler

import (
	"net/http"

	hls_settings "github.com/number571/go-peer/cmd/hls/internal/settings"
	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response(w, pkg_settings.CErrorNone, hls_settings.CTitlePattern)
	}
}
