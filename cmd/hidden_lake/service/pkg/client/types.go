package client

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

type IClient interface {
	GetIndex() (string, error)

	GetPubKey() (asymmetric.IPubKey, asymmetric.IEphPubKey, error)
	SetPrivKey(asymmetric.IPrivKey, asymmetric.IEphPubKey) error

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	GetConnections() ([]string, error)
	AddConnection(string) error
	DelConnection(string) error

	HandleMessage(message.IMessage) error

	BroadcastRequest(asymmetric.IPubKey, request.IRequest) error
	FetchRequest(asymmetric.IPubKey, request.IRequest) (response.IResponse, error)
}

type IRequester interface {
	GetIndex() (string, error)

	GetPubKey() (asymmetric.IPubKey, asymmetric.IEphPubKey, error)
	SetPrivKey(*pkg_settings.SPrivKey) error

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(*pkg_settings.SFriend) error
	DelFriend(*pkg_settings.SFriend) error

	GetConnections() ([]string, error)
	AddConnection(string) error
	DelConnection(string) error

	HandleMessage(string) error

	BroadcastRequest(*pkg_settings.SRequest) error
	FetchRequest(*pkg_settings.SRequest) (response.IResponse, error)
}

type IBuilder interface {
	SetPrivKey(asymmetric.IPrivKey, asymmetric.IEphPubKey) *pkg_settings.SPrivKey
	Friend(string, asymmetric.IPubKey) *pkg_settings.SFriend
	Message(message.IMessage) string
	Request(asymmetric.IPubKey, request.IRequest) *pkg_settings.SRequest
}
