package api

import (
	"encoding/json"
	"net/http"
)

func Response(pW http.ResponseWriter, pRet int, pRes string) {
	pW.Header().Set("Content-Type", CContentType)
	json.NewEncoder(pW).Encode(&SResponse{
		FResult: pRes,
		FReturn: pRet,
	})
}
