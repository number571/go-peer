package api

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", CContentType)
	json.NewEncoder(w).Encode(&SResponse{
		FResult: res,
		FReturn: ret,
	})
}
