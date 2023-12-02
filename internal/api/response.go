package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
)

func Response(pW http.ResponseWriter, pRet int, pRes interface{}) {
	var (
		contentType = ""
		respBytes   []byte
	)

	switch x := pRes.(type) {
	case []byte:
		contentType = cTextPlain
		respBytes = x
	case string:
		contentType = cTextPlain
		respBytes = []byte(x)
	default:
		contentType = cApplicationJSON
		respBytes = encoding.SerializeJSON(x)
	}

	pW.Header().Set("Content-Type", contentType)
	pW.WriteHeader(pRet)
	io.Copy(pW, bytes.NewBuffer(respBytes))
}
