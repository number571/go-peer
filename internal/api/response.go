package api

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
)

func Response(pW http.ResponseWriter, pRet int, pRes interface{}) {
	var (
		contentType = ""
		respStr     = ""
	)

	switch x := pRes.(type) {
	case []byte:
		contentType = cTextPlain
		respStr = encoding.HexEncode(x)
	case string:
		contentType = cTextPlain
		respStr = x
	default:
		contentType = cApplicationJSON
		respStr = string(encoding.Serialize(x, false))
	}

	pW.Header().Set("Content-Type", contentType)
	pW.WriteHeader(pRet)
	fmt.Fprint(pW, respStr)
}
