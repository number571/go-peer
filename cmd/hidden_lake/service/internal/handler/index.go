package handler

import (
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.Response(w, pkg_settings.CErrorNone, pkg_settings.CTitlePattern)
	}
}
