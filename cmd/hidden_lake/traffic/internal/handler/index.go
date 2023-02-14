package handler

import (
	"net/http"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.Response(w, hlt_settings.CErrorNone, hlt_settings.CTitlePattern)
	}
}
