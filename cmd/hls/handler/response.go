package handler

import (
	"encoding/json"
	"net/http"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", hls_settings.CContentType)
	json.NewEncoder(w).Encode(&hls_settings.SResponse{
		FResult: res,
		FReturn: ret,
	})
}
