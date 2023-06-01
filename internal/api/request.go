package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/errors"
)

func Request(pClient *http.Client, pMethod, pURL string, pData interface{}) (string, error) {
	var requestBytes []byte

	switch x := pData.(type) {
	case []byte:
		requestBytes = x
	case string:
		requestBytes = []byte(x)
	default:
		jsonValue, err := json.Marshal(pData)
		if err != nil {
			return "", errors.WrapError(err, "marshal request")
		}
		requestBytes = jsonValue
	}

	req, err := http.NewRequest(
		pMethod,
		pURL,
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return "", errors.WrapError(err, "new request")
	}

	req.Header.Set("Content-Type", CContentType)
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
