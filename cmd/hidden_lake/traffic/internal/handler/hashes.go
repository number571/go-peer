package handler

import (
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/internal/database"
	"github.com/number571/go-peer/internal/api"
)

func HandleHashesAPI(pWrapperDB database.IWrapperDB) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet {
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		database := pWrapperDB.Get()
		hashes, err := database.Hashes()
		if err != nil {
			api.Response(pW, http.StatusInternalServerError, "failed: load size from DB")
			return
		}

		api.Response(pW, http.StatusOK, strings.Join(hashes, ";"))
	}
}
