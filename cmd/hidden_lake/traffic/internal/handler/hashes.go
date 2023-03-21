package handler

import (
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/api"
)

func HandleHashesAPI(pWrapperDB database.IWrapperDB) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet {
			api.Response(pW, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()
		hashes, err := database.Hashes()
		if err != nil {
			api.Response(pW, pkg_settings.CErrorLoad, "failed: load size")
			return
		}

		api.Response(pW, pkg_settings.CErrorNone, strings.Join(hashes, ";"))
	}
}
