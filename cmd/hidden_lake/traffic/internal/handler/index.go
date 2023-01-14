package handler

import (
	"net/http"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response(w, hlt_settings.CErrorNone, hlt_settings.CTitlePattern)
	}
}
