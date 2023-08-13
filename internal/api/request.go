package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
)

func Request(pClient *http.Client, pMethod, pURL string, pData interface{}) (string, error) {
	var (
		contentType = ""
		reqBytes    []byte
	)

	switch x := pData.(type) {
	case []byte:
		contentType = cTextPlain
		reqBytes = x
	case string:
		contentType = cTextPlain
		reqBytes = []byte(x)
	default:
		contentType = cApplicationJSON
		reqBytes = encoding.Serialize(x, false)
	}

	req, err := http.NewRequest(
		pMethod,
		pURL,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return "", errors.WrapError(err, "new request")
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := pClient.Do(req)
	if err != nil {
		return "", errors.WrapError(err, "do request")
	}
	defer resp.Body.Close()

	result, err := loadResponse(resp.StatusCode, resp.Body)
	if err != nil {
		return "", errors.WrapError(err, "load response")
	}
	return result, nil
}

func loadResponse(pStatusCode int, pReader io.ReadCloser) (string, error) {
	resp, err := io.ReadAll(pReader)
	if err != nil {
		return "", errors.WrapError(err, "read response")
	}

	result := string(resp)
	if pStatusCode != http.StatusOK {
		return "", errors.NewError(fmt.Sprintf("error code = %d (%s)", pStatusCode, result))
	}

	return result, nil
}
