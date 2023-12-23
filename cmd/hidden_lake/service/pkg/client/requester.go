package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate          = "%s" + hls_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "%s" + hls_settings.CHandleConfigSettingsPath
	cHandleConfigConnectsTemplate = "%s" + hls_settings.CHandleConfigConnectsPath
	cHandleConfigFriendsTemplate  = "%s" + hls_settings.CHandleConfigFriendsPath
	cHandleNetworkOnlineTemplate  = "%s" + hls_settings.CHandleNetworkOnlinePath
	cHandleNetworkRequestTemplate = "%s" + hls_settings.CHandleNetworkRequestPath
	cHandleNetworkPubKeyTemplate  = "%s" + hls_settings.CHandleNetworkPubKeyPath
)

type sRequester struct {
	fHost   string
	fClient *http.Client
}

func NewRequester(pHost string, pClient *http.Client) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
	}
}

func (p *sRequester) GetIndex() (string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("get index (requester): %w", err)
	}

	result := string(res)
	if result != hls_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
	}

	return result, nil
}

func (p *sRequester) GetSettings() (config.IConfigSettings, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get settings (requester): %w", err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON([]byte(res), cfgSettings); err != nil {
		return nil, fmt.Errorf("decode settings (requester): %w", err)
	}

	return cfgSettings, nil
}

func (p *sRequester) SetNetworkKey(pNetworkKey string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		pNetworkKey,
	)
	if err != nil {
		return fmt.Errorf("set network key (requester): %w", err)
	}
	return nil
}

func (p *sRequester) FetchRequest(pRequest *hls_settings.SRequest) (response.IResponse, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch request (requester): %w", err)
	}

	resp, err := response.LoadResponse([]byte(res))
	if err != nil {
		return nil, fmt.Errorf("load fetch response (requester): %w", err)
	}
	return resp, nil
}

func (p *sRequester) BroadcastRequest(pRequest *hls_settings.SRequest) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return fmt.Errorf("broadcast request (requester): %w", err)
	}
	return nil
}

func (p *sRequester) GetFriends() (map[string]asymmetric.IPubKey, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get friends (requester): %w", err)
	}

	var vFriends []hls_settings.SFriend
	if err := encoding.DeserializeJSON([]byte(res), &vFriends); err != nil {
		return nil, fmt.Errorf("deserialize friends (requeser): %w", err)
	}

	result := make(map[string]asymmetric.IPubKey, len(vFriends))
	for _, friend := range vFriends {
		result[friend.FAliasName] = asymmetric.LoadRSAPubKey(friend.FPublicKey)
	}

	return result, nil
}

func (p *sRequester) AddFriend(pFriend *hls_settings.SFriend) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return fmt.Errorf("add friend (requester): %w", err)
	}
	return nil
}

func (p *sRequester) DelFriend(pFriend *hls_settings.SFriend) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return fmt.Errorf("del friend (requester): %w", err)
	}
	return nil
}

func (p *sRequester) GetOnlines() ([]string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get onlines (requester): %w", err)
	}

	var onlines []string
	if err := encoding.DeserializeJSON([]byte(res), &onlines); err != nil {
		return nil, fmt.Errorf("deserialize onlines (requester): %w", err)
	}

	return onlines, nil
}

func (p *sRequester) DelOnline(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return fmt.Errorf("del online (requester): %w", err)
	}
	return nil
}

func (p *sRequester) GetConnections() ([]string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get connections (requester): %w", err)
	}

	var connects []string
	if err := encoding.DeserializeJSON([]byte(res), &connects); err != nil {
		return nil, fmt.Errorf("deserialize connections (requeser): %w", err)
	}

	return connects, nil
}

func (p *sRequester) AddConnection(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return fmt.Errorf("add connection (requester): %w", err)
	}
	return nil
}

func (p *sRequester) DelConnection(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return fmt.Errorf("del connection (requester): %w", err)
	}
	return nil
}

func (p *sRequester) GetPubKey() (asymmetric.IPubKey, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkPubKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get public key (requester): %w", err)
	}

	pubKey := asymmetric.LoadRSAPubKey(string(res))
	if pubKey == nil {
		return nil, errors.New("got invalid public key")
	}

	return pubKey, nil
}
