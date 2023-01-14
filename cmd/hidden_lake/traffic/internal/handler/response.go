package handler

import (
	"encoding/json"
	"net/http"

	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
)

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", pkg_settings.CContentType)
	json.NewEncoder(w).Encode(&pkg_settings.SResponse{
		FResult: res,
		FReturn: ret,
	})
}
