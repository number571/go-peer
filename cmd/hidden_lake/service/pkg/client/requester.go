package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
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

func (p *sRequester) GetIndex(pCtx context.Context) (string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", utils.MergeErrors(ErrBadRequest, err)
	}

	result := string(res)
	if result != hls_settings.CServiceFullName {
		return "", ErrInvalidTitle
	}

	return result, nil
}

func (p *sRequester) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON([]byte(res), cfgSettings); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return cfgSettings, nil
}

func (p *sRequester) SetNetworkKey(pCtx context.Context, pNetworkKey string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		pNetworkKey,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) FetchRequest(pCtx context.Context, pRequest *hls_settings.SRequest) (response.IResponse, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	resp, err := response.LoadResponse([]byte(res))
	if err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}
	return resp, nil
}

func (p *sRequester) BroadcastRequest(pCtx context.Context, pRequest *hls_settings.SRequest) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) GetFriends(pCtx context.Context) (map[string]asymmetric.IPubKey, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	var vFriends []hls_settings.SFriend
	if err := encoding.DeserializeJSON([]byte(res), &vFriends); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	result := make(map[string]asymmetric.IPubKey, len(vFriends))
	for _, friend := range vFriends {
		result[friend.FAliasName] = asymmetric.LoadRSAPubKey(friend.FPublicKey)
	}

	return result, nil
}

func (p *sRequester) AddFriend(pCtx context.Context, pFriend *hls_settings.SFriend) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) DelFriend(pCtx context.Context, pFriend *hls_settings.SFriend) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) GetOnlines(pCtx context.Context) ([]string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	var onlines []string
	if err := encoding.DeserializeJSON([]byte(res), &onlines); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return onlines, nil
}

func (p *sRequester) DelOnline(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) GetConnections(pCtx context.Context) ([]string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	var connects []string
	if err := encoding.DeserializeJSON([]byte(res), &connects); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return connects, nil
}

func (p *sRequester) AddConnection(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) DelConnection(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) GetPubKey(pCtx context.Context) (asymmetric.IPubKey, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkPubKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	pubKey := asymmetric.LoadRSAPubKey(string(res))
	if pubKey == nil {
		return nil, ErrInvalidPublicKey
	}

	return pubKey, nil
}
