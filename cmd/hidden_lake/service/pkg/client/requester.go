package client

import (
	"fmt"
	"net/http"
	"strings"

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
		return "", err
	}

	if res != pkg_settings.CTitlePattern {
		return "", fmt.Errorf("incorrect title pattern")
	}
	return res, nil
}

func (p *sRequester) FetchRequest(pRequest *pkg_settings.SRequest) ([]byte, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return nil, err
	}

	return encoding.HexDecode(res), nil
}

func (p *sRequester) BroadcastRequest(pRequest *pkg_settings.SRequest) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(pkg_settings.CHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	return err
}

func (p *sRequester) GetFriends() (map[string]asymmetric.IPubKey, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}

	listFriends := deleteVoidStrings(strings.Split(res, ","))
	result := make(map[string]asymmetric.IPubKey, len(listFriends))
	for _, friend := range listFriends {
		splited := strings.Split(friend, ":")
		if len(splited) != 2 {
			return nil, fmt.Errorf("length of splited != 2")
		}
		aliasName := splited[0]
		pubKeyStr := splited[1]
		result[aliasName] = asymmetric.LoadRSAPubKey(pubKeyStr)
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
	return err
}

func (p *sRequester) DelFriend(pFriend *pkg_settings.SFriend) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	return err
}

func (p *sRequester) GetOnlines() ([]string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return deleteVoidStrings(strings.Split(res, ",")), nil
}

func (p *sRequester) DelOnline(pConnect *pkg_settings.SConnect) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, p.fHost),
		pConnect,
	)
	return err
}

func (p *sRequester) GetConnections() ([]string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return deleteVoidStrings(strings.Split(res, ",")), nil
}

func (p *sRequester) AddConnection(pConnect *pkg_settings.SConnect) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	return err
}

func (p *sRequester) DelConnection(pConnect *pkg_settings.SConnect) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	return err
}

func (p *sRequester) SetPrivKey(pPrivKey *pkg_settings.SPrivKey) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, p.fHost),
		pPrivKey,
	)
	return err
}

func (p *sRequester) GetPubKey() (asymmetric.IPubKey, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return asymmetric.LoadRSAPubKey(res), nil
}

func deleteVoidStrings(pS []string) []string {
	result := make([]string, 0, len(pS))
	for _, v := range pS {
		r := strings.TrimSpace(v)
		if r == "" {
			continue
		}
		result = append(result, r)
	}
	return result
}
