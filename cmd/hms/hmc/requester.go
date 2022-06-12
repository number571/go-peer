package hmc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/offline/message"
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

func (r *sRequester) Size(request *hms_settings.SSizeRequest) (uint64, error) {
	resp, err := http.Post(
		fmt.Sprintf(hms_settings.CSizeTemplate, r.host),
		hms_settings.CContentType,
		bytes.NewReader(encoding.Serialize(request)),
	)
	if err != nil {
		return 0, err
	}

	var response hms_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.Return != hms_settings.CErrorNone {
		return 0, fmt.Errorf("%s", string(response.Result))
	}

	return encoding.BytesToUint64(response.Result), nil
}

func (r *sRequester) Load(request *hms_settings.SLoadRequest) (message.IMessage, error) {
	resp, err := http.Post(
		fmt.Sprintf(hms_settings.CLoadTemplate, r.host),
		hms_settings.CContentType,
		bytes.NewReader(encoding.Serialize(request)),
	)
	if err != nil {
		return nil, err
	}

	var response hms_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.Return != hms_settings.CErrorNone {
		return nil, fmt.Errorf("%s", string(response.Result))
	}

	msg := message.LoadPackage(response.Result).ToMessage()
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	return msg, nil
}

func (r *sRequester) Push(request *hms_settings.SPushRequest) error {
	resp, err := http.Post(
		fmt.Sprintf(hms_settings.CPushTemplate, r.host),
		hms_settings.CContentType,
		bytes.NewReader(encoding.Serialize(request)),
	)
	if err != nil {
		return err
	}
	var response hms_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.Return != hms_settings.CErrorNone {
		return fmt.Errorf("%s", string(response.Result))
	}

	return nil
}
