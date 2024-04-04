package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

func Response(pW http.ResponseWriter, pRet int, pRes interface{}) error {
	var (
		contentType string
		respBytes   []byte
	)

	switch x := pRes.(type) {
	case []byte:
		contentType = CApplicationOctetStream
		respBytes = x
	case string:
		contentType = CTextPlain
		respBytes = []byte(x)
	default:
		contentType = CApplicationJSON
		respBytes = encoding.SerializeJSON(x)
	}

	pW.Header().Set("Content-Type", contentType)
	pW.WriteHeader(pRet)

	_, err := io.Copy(pW, bytes.NewBuffer(respBytes))
	if err != nil {
		return utils.MergeErrors(ErrCopyBytes, err)
	}

	return nil
}
