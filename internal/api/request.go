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

	res, err := loadResponse(resp.Body)
	if err != nil {
		return "", err
	}

	return res.FResult, nil
}

func loadResponse(pReader io.ReadCloser) (*SResponse, error) {
	body, err := io.ReadAll(pReader)
	if err != nil {
		return nil, err
	}

	resp := &SResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	if resp.FReturn != CErrorNone {
		return nil, fmt.Errorf("error code = %d (%s)", resp.FReturn, resp.FResult)
	}

	return resp, nil
}
