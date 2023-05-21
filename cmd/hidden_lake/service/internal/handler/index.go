package handler

import (
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

func HandleIndexAPI() http.HandlerFunc {
	return func(pW http.ResponseWriter, _ *http.Request) {
		api.Response(pW, http.StatusOK, pkg_settings.CTitlePattern)
	}
}
