package handler

import (
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func HashesPage(db database.IKeyValueDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			response(w, pkg_settings.CErrorMethod, "failed: incorrect method")
			return
		}

		hashes, err := db.Hashes()
		if err != nil {
			response(w, pkg_settings.CErrorLoad, "failed: load size")
			return
		}

		response(w, pkg_settings.CErrorNone, strings.Join(hashes, ","))
	}
}
