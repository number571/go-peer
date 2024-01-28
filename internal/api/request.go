package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
)

func Request(pClient *http.Client, pMethod, pURL string, pData interface{}) ([]byte, error) {
	var (
		contentType = ""
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

	req, err := http.NewRequest(
		pMethod,
		pURL,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := pClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	result, err := loadResponse(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("load request: %w", err)
	}
	return result, nil
}

func loadResponse(pStatusCode int, pReader io.ReadCloser) ([]byte, error) {
	resp, err := io.ReadAll(pReader)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if pStatusCode != http.StatusOK {
		return nil, fmt.Errorf("error code = %d (%x)", pStatusCode, resp)
	}

	return resp, nil
}
