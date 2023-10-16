package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

type IClient interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	GetNetworkKey() (string, error)
	SetNetworkKey(string) error

	GetPubKey() (asymmetric.IPubKey, error)

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	GetConnections(bool) ([]string, error)
	AddConnection(bool, string) error
	DelConnection(bool, string) error

	HandleMessage(message.IMessage) error

	BroadcastRequest(string, request.IRequest) error
	FetchRequest(string, request.IRequest) (response.IResponse, error)
}

type IRequester interface {
	GetIndex() (string, error)
	GetSettings() (config.IConfigSettings, error)

	GetNetworkKey() (string, error)
	SetNetworkKey(string) error

	GetPubKey() (asymmetric.IPubKey, error)

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(*pkg_settings.SFriend) error
	DelFriend(*pkg_settings.SFriend) error

	GetConnections(bool) ([]string, error)
	AddConnection(bool, string) error
	DelConnection(bool, string) error

	HandleMessage(string) error

	BroadcastRequest(*pkg_settings.SRequest) error
	FetchRequest(*pkg_settings.SRequest) (response.IResponse, error)
}

type IBuilder interface {
	Request(string, request.IRequest) *pkg_settings.SRequest
	Friend(string, asymmetric.IPubKey) *pkg_settings.SFriend
	Message(message.IMessage) string
}
