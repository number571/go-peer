package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
			return "", err
		}
		requestBytes = jsonValue
	}

	req, err := http.NewRequest(
		pMethod,
		pURL,
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", CContentType)
	resp, err := pClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return loadResponse(resp.StatusCode, resp.Body)
}

func loadResponse(pStatusCode int, pReader io.ReadCloser) (string, error) {
	resp, err := io.ReadAll(pReader)
	if err != nil {
		return "", err
	}

	result := string(resp)
	if pStatusCode != http.StatusOK {
		return "", fmt.Errorf("error code = %d (%s)", pStatusCode, result)
	}

	return result, nil
}
