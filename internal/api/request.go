package api

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
)

func Request(
	pCtx context.Context,
	pClient *http.Client,
	pMethod, pURL string,
	pData interface{},
) ([]byte, error) {
	var (
		contentType string
		reqBytes    []byte
	)

	switch x := pData.(type) {
	case []byte:
		contentType = CTextPlain
		reqBytes = x
	case string:
		contentType = CTextPlain
		reqBytes = []byte(x)
	default:
		contentType = CApplicationJSON
		reqBytes = encoding.SerializeJSON(x)
	}

	req, err := http.NewRequestWithContext(
		pCtx,
		pMethod,
		pURL,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBuildRequest, err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := pClient.Do(req)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}
	defer resp.Body.Close()

	result, err := loadResponse(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, utils.MergeErrors(ErrLoadResponse, err)
	}
	return result, nil
}

func loadResponse(pStatusCode int, pReader io.ReadCloser) ([]byte, error) {
	resp, err := io.ReadAll(pReader)
	if err != nil {
		return nil, utils.MergeErrors(ErrReadResponse, err)
	}

	if pStatusCode != http.StatusOK {
		return nil, ErrBadStatusCode
	}

	return resp, nil
}
