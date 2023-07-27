package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
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
		return "", errors.WrapError(err, "get index (requester)")
	}

	if res != pkg_settings.CTitlePattern {
		return "", errors.NewError("incorrect title pattern")
	}
	return res, nil
}

func (p *sRequester) HandleMessage(pMsg string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNetworkMessageTemplate, p.fHost),
		pMsg,
	)
	if err != nil {
		return errors.WrapError(err, "handle message (requester)")
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
		return nil, errors.WrapError(err, "fetch request (requester)")
	}

	resp, err := response.LoadResponse([]byte(res))
	if err != nil {
		return nil, errors.WrapError(err, "load fetch response (requester)")
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
		return errors.WrapError(err, "broadcast request (requester)")
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
		return nil, errors.WrapError(err, "get friends (requester)")
	}

	listFriends := deleteVoidStrings(strings.Split(res, ","))
	result := make(map[string]asymmetric.IPubKey, len(listFriends))
	for _, friend := range listFriends {
		splited := strings.Split(friend, ":")
		if len(splited) != 2 {
			return nil, errors.NewError("length of splited != 2")
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
	if err != nil {
		return errors.WrapError(err, "add friend (requester)")
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
		return errors.WrapError(err, "del friend (requester)")
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
		return nil, errors.WrapError(err, "get onlines (requester)")
	}
	return deleteVoidStrings(strings.Split(res, ",")), nil
}

func (p *sRequester) DelOnline(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleNetworkOnlineTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return errors.WrapError(err, "del online (requester)")
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
		return nil, errors.WrapError(err, "get connections (requester)")
	}
	return deleteVoidStrings(strings.Split(res, ",")), nil
}

func (p *sRequester) AddConnection(pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return errors.WrapError(err, "add connection (requester)")
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
		return errors.WrapError(err, "del connection (requester)")
	}
	return nil
}

func (p *sRequester) SetPrivKey(pPrivKey *pkg_settings.SPrivKey) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, p.fHost),
		pPrivKey,
	)
	if err != nil {
		return errors.WrapError(err, "set private key (requester)")
	}
	return nil
}

func (p *sRequester) GetPubKey() (asymmetric.IPubKey, asymmetric.IPubKey, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleNodeKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, nil, errors.WrapError(err, "get public key (requester)")
	}

	splited := strings.Split(res, ",")
	if len(splited) != 2 {
		return nil, nil, errors.NewError("got length of splited != 2")
	}

	pubKey := asymmetric.LoadRSAPubKey(splited[0])
	if pubKey == nil {
		return nil, nil, errors.NewError("got invalid public key")
	}

	ephPubKey := asymmetric.LoadRSAPubKey(splited[1])
	if ephPubKey == nil {
		return nil, nil, errors.NewError("got invalid eph public key")
	}

	return pubKey, ephPubKey, nil
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
