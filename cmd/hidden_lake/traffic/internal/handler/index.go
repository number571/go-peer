package handler

import (
	"net/http"

	hlt_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(pW http.ResponseWriter, _ *http.Request) {
		api.Response(pW, http.StatusOK, hlt_settings.CTitlePattern)
	}
}
