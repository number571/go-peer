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

func NewRequester(host string, client *http.Client) IRequester {
	return &sRequester{
		fHost:   host,
		fClient: client,
	}
}

func (requester *sRequester) GetIndex() (string, error) {
	res, err := api.Request(
		requester.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleIndexTemplate, requester.fHost),
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

func (requester *sRequester) FetchRequest(push *pkg_settings.SRequest) ([]byte, error) {
	res, err := api.Request(
		requester.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNetworkRequestTemplate, requester.fHost),
		push,
	)
	if err != nil {
		return nil, err
	}

	return encoding.HexDecode(res), nil
}

func (requester *sRequester) BroadcastRequest(push *pkg_settings.SRequest) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodPut,
		fmt.Sprintf(pkg_settings.CHandleNetworkRequestTemplate, requester.fHost),
		push,
	)
	return err
}

func (requester *sRequester) GetFriends() (map[string]asymmetric.IPubKey, error) {
	res, err := api.Request(
		requester.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, requester.fHost),
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

func (requester *sRequester) AddFriend(friend *pkg_settings.SFriend) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, requester.fHost),
		friend,
	)
	return err
}

func (requester *sRequester) DelFriend(friend *pkg_settings.SFriend) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleConfigFriendsTemplate, requester.fHost),
		friend,
	)
	return err
}

func (requester *sRequester) GetOnlines() ([]string, error) {
	res, err := api.Request(
		requester.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, requester.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return deleteVoidStrings(strings.Split(res, ",")), nil
}

func (requester *sRequester) DelOnline(connect *pkg_settings.SConnect) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, requester.fHost),
		connect,
	)
	return err
}

func (requester *sRequester) GetConnections() ([]string, error) {
	res, err := api.Request(
		requester.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, requester.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return deleteVoidStrings(strings.Split(res, ",")), nil
}

func (requester *sRequester) AddConnection(connect *pkg_settings.SConnect) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, requester.fHost),
		connect,
	)
	return err
}

func (requester *sRequester) DelConnection(connect *pkg_settings.SConnect) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, requester.fHost),
		connect,
	)
	return err
}

func (requester *sRequester) SetPrivKey(privKey *pkg_settings.SPrivKey) error {
	_, err := api.Request(
		requester.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, requester.fHost),
		privKey,
	)
	return err
}

func (requester *sRequester) GetPubKey() (asymmetric.IPubKey, error) {
	res, err := api.Request(
		requester.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, requester.fHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return asymmetric.LoadRSAPubKey(res), nil
}

func deleteVoidStrings(s []string) []string {
	result := make([]string, 0, len(s))
	for _, v := range s {
		r := strings.TrimSpace(v)
		if r == "" {
			continue
		}
		result = append(result, r)
	}
	return result
}
