package hmc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hms/utils"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	host   string
	client local.IClient
}

func NewClient(client local.IClient, host string) IClient {
	return &sClient{
		host:   host,
		client: client,
	}
}

func (client *sClient) Size() (uint64, error) {
	pubBytes := client.client.PubKey().Bytes()
	hashRecv := crypto.NewHasher(pubBytes).Bytes()

	request := struct {
		Receiver []byte `json:"receiver"`
	}{
		Receiver: hashRecv,
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/size", client.host),
		"application/json",
		bytes.NewReader(utils.Serialize(request)),
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

func (client *sClient) Load(n uint64) ([]byte, error) {
	pubBytes := client.client.PubKey().Bytes()
	hashRecv := crypto.NewHasher(pubBytes).Bytes()

	request := hms_settings.SLoadRequest{
		Receiver: hashRecv,
		Index:    n,
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/load", client.host),
		"application/json",
		bytes.NewReader(utils.Serialize(request)),
	)
	if err != nil {
		return nil, err
	}

	var response hms_settings.SResponse
	json.NewDecoder(resp.Body).Decode(&response)

	if response.Return != hms_settings.CErrorNone {
		return nil, fmt.Errorf("%s", string(response.Result))
	}

	msg := local.LoadPackage(response.Result).ToMessage()
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	msg, title := client.client.Decrypt(msg)
	if string(title) != hms_settings.CPatternTitle {
		return nil, fmt.Errorf("title is not equal")
	}

	return msg.Body().Data(), nil
}

func (client *sClient) Push(receiver crypto.IPubKey, msg []byte) error {
	pubBytes := receiver.Bytes()
	hashRecv := crypto.NewHasher(pubBytes).Bytes()

	encMsg, _ := client.client.Encrypt(
		local.NewRoute(receiver),
		local.NewMessage([]byte(hms_settings.CPatternTitle), msg),
	)

	request := hms_settings.SPushRequest{
		Receiver: hashRecv,
		Package:  encMsg.ToPackage().Bytes(),
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/push", client.host),
		"application/json",
		bytes.NewReader(utils.Serialize(request)),
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
