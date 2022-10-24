package hlc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	host string
}

func NewRequester(host string) IRequester {
	return &sRequester{
		host: host,
	}
}

func (requester *sRequester) Request(req *hls_settings.SPush) ([]byte, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	respPost, err := http.Post(
		requester.host+hls_settings.CHandlePush,
		hls_settings.CContentType,
		bytes.NewBuffer(jsonValue),
	)
	if err != nil {
		return nil, err
	}
	defer respPost.Body.Close()

	resp, err := loadResponse(respPost.Body)
	if err != nil {
		return nil, err
	}

	return encoding.HexDecode(resp.FResult), nil
}

func (requester *sRequester) Friends() ([]asymmetric.IPubKey, error) {
	respGet, err := http.Get(requester.host + hls_settings.CHandleFriends)
	if err != nil {
		return nil, err
	}
	defer respGet.Body.Close()

	resp, err := loadResponse(respGet.Body)
	if err != nil {
		return nil, err
	}

	listPubKeysStr := strings.Split(resp.FResult, ",")
	listPubKeys := make([]asymmetric.IPubKey, 0, len(listPubKeysStr))
	for _, pubKeyStr := range listPubKeysStr {
		listPubKeys = append(listPubKeys, asymmetric.LoadRSAPubKey(pubKeyStr))
	}

	return listPubKeys, nil
}

func (requester *sRequester) Online() ([]string, error) {
	respGet, err := http.Get(requester.host + hls_settings.CHandleOnline)
	if err != nil {
		return nil, err
	}
	defer respGet.Body.Close()

	resp, err := loadResponse(respGet.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(resp.FResult, ","), nil
}

func (requester *sRequester) PubKey() (asymmetric.IPubKey, error) {
	respGet, err := http.Get(requester.host + hls_settings.CHandlePubKey)
	if err != nil {
		return nil, err
	}
	defer respGet.Body.Close()

	resp, err := loadResponse(respGet.Body)
	if err != nil {
		return nil, err
	}

	return asymmetric.LoadRSAPubKey(resp.FResult), nil
}

func loadResponse(reader io.ReadCloser) (*hls_settings.SResponse, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	resp := &hls_settings.SResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	if resp.FReturn != hls_settings.CErrorNone {
		return nil, fmt.Errorf("error code = %d", resp.FReturn)
	}

	return resp, nil
}
