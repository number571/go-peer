package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IRequester = &sRequester{}
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
		fmt.Sprintf(pkg_settings.CHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("get index (requester): %w", err)
	}

	result := string(res)
	if result != pkg_settings.CTitlePattern {
		return "", errors.New("incorrect title pattern")
	}

	return result, nil
}

func (p *sRequester) GetSettings() (config.IConfigSettings, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get settings (requester): %w", err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.Deserialize([]byte(res), cfgSettings); err != nil {
		return nil, fmt.Errorf("decode settings (requester): %w", err)
	}

	return cfgSettings, nil
}

func (p *sRequester) GetNetworkKey() (string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleNetworkKeyTemplate, p.fHost),
		nil,
	)

	result := string(res)
	if err != nil {
		return "", fmt.Errorf("get network key (requester): %w", err)
	}

	return result, nil
}

func (p *sRequester) SetNetworkKey(pNetworkKey string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNetworkKeyTemplate, p.fHost),
		pNetworkKey,
	)
	if err != nil {
		return fmt.Errorf("set network key (requester): %w", err)
	}
	return nil
}

func (p *sRequester) FetchRequest(pRequest *pkg_settings.SRequest) (response.IResponse, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNetworkRequestTemplate, p.fHost),
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

func (p *sRequester) BroadcastRequest(pRequest *pkg_settings.SRequest) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(pkg_settings.CHandleNetworkRequestTemplate, p.fHost),
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
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get friends (requester): %w", err)
	}

	var vFriends []pkg_settings.SFriend
	if err := encoding.Deserialize([]byte(res), &vFriends); err != nil {
		return nil, fmt.Errorf("deserialize friends (requeser): %w", err)
	}

	result := make(map[string]asymmetric.IPubKey, len(vFriends))
	for _, friend := range vFriends {
		result[friend.FAliasName] = asymmetric.LoadRSAPubKey(friend.FPublicKey)
	}

	return result, nil
}

func (p *sRequester) AddFriend(pFriend *pkg_settings.SFriend) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return fmt.Errorf("add friend (requester): %w", err)
	}
	return nil
}

func (p *sRequester) DelFriend(pFriend *pkg_settings.SFriend) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, p.fHost),
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
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get onlines (requester): %w", err)
	}

	var onlines []string
	if err := encoding.Deserialize([]byte(res), &onlines); err != nil {
		return nil, fmt.Errorf("deserialize onlines (requester): %w", err)
	}

	return onlines, nil
}

func (p *sRequester) DelOnline(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, p.fHost),
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
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get connections (requester): %w", err)
	}

	var connects []string
	if err := encoding.Deserialize([]byte(res), &connects); err != nil {
		return nil, fmt.Errorf("deserialize connections (requeser): %w", err)
	}

	return connects, nil
}

func (p *sRequester) AddConnection(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
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
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
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
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, p.fHost),
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
