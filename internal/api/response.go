package api

import (
	"fmt"
	"net/http"
)

func Response(pW http.ResponseWriter, pRet int, pRes string) {
	pW.WriteHeader(pRet)
	pW.Header().Set("Content-Type", CContentType)
	fmt.Fprint(pW, pRes)
}
