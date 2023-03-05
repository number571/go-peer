package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Request(method, url string, data interface{}) (string, error) {
	var requestBytes []byte

	switch x := data.(type) {
	case []byte:
		requestBytes = x
	default:
		jsonValue, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		requestBytes = jsonValue
	}

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", CContentType)
	resp, err := http.DefaultClient.Do(req)
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

func loadResponse(reader io.ReadCloser) (*SResponse, error) {
	body, err := io.ReadAll(reader)
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
