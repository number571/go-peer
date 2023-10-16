package client

import (
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
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

func (p *sRequester) GetSettings() (config.IConfigSettings, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.WrapError(err, "get settings (requester)")
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.Deserialize([]byte(res), cfgSettings); err != nil {
		return nil, errors.WrapError(err, "decode settings (requester)")
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
	if err != nil {
		return "", errors.WrapError(err, "get network key (requester)")
	}
	return res, nil
}

func (p *sRequester) SetNetworkKey(pNetworkKey string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleNetworkKeyTemplate, p.fHost),
		pNetworkKey,
	)
	if err != nil {
		return errors.WrapError(err, "set network key (requester)")
	}
	return nil
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

	var vFriends []pkg_settings.SFriend
	if err := encoding.Deserialize([]byte(res), &vFriends); err != nil {
		return nil, errors.WrapError(err, "deserialize friends (requeser)")
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

	var onlines []string
	if err := encoding.Deserialize([]byte(res), &onlines); err != nil {
		return nil, errors.WrapError(err, "deserialize onlines (requeser)")
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
		return errors.WrapError(err, "del online (requester)")
	}
	return nil
}

func (p *sRequester) GetConnections(pIsBackup bool) ([]string, error) {
	res, err := api.Request(
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate+"?is_backup=%s", p.fHost, isBackup(pIsBackup)),
		nil,
	)
	if err != nil {
		return nil, errors.WrapError(err, "get connections (requester)")
	}

	var connects []string
	if err := encoding.Deserialize([]byte(res), &connects); err != nil {
		return nil, errors.WrapError(err, "deserialize connections (requeser)")
	}

	return connects, nil
}

func (p *sRequester) AddConnection(pIsBackup bool, pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate+"?is_backup=%s", p.fHost, isBackup(pIsBackup)),
		pConnect,
	)
	if err != nil {
		return errors.WrapError(err, "add connection (requester)")
	}
	return nil
}

func (p *sRequester) DelConnection(pIsBackup bool, pConnect string) error {
	_, err := api.Request(
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(pkg_settings.CHandleConfigConnectsTemplate+"?is_backup=%s", p.fHost, isBackup(pIsBackup)),
		pConnect,
	)
	if err != nil {
		return errors.WrapError(err, "del connection (requester)")
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
		return nil, errors.WrapError(err, "get public key (requester)")
	}

	pubKey := asymmetric.LoadRSAPubKey(res)
	if pubKey == nil {
		return nil, errors.NewError("got invalid public key")
	}

	return pubKey, nil
}

func isBackup(pIsBackup bool) string {
	if pIsBackup {
		return "true"
	}
	return "false"
}
