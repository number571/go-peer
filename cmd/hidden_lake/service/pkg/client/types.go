package client

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

type IClient interface {
	PubKey() (asymmetric.IPubKey, error)
	PrivKey(asymmetric.IPrivKey) error

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	GetConnections() ([]string, error)
	AddConnection(string) error
	DelConnection(string) error

	Broadcast(asymmetric.IPubKey, request.IRequest) error
	Request(asymmetric.IPubKey, request.IRequest) ([]byte, error)
}

type IRequester interface {
	PubKey() (asymmetric.IPubKey, error)
	PrivKey(*pkg_settings.SPrivKey) error

	GetOnlines() ([]string, error)
	DelOnline(*pkg_settings.SConnect) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(*pkg_settings.SFriend) error
	DelFriend(*pkg_settings.SFriend) error

	GetConnections() ([]string, error)
	AddConnection(*pkg_settings.SConnect) error
	DelConnection(*pkg_settings.SConnect) error

	Broadcast(*pkg_settings.SPush) error
	Request(*pkg_settings.SPush) ([]byte, error)
}

type IBuilder interface {
	PrivKey(asymmetric.IPrivKey) *pkg_settings.SPrivKey
	Connect(string) *pkg_settings.SConnect
	Friend(string, asymmetric.IPubKey) *pkg_settings.SFriend
	Push(asymmetric.IPubKey, request.IRequest) *pkg_settings.SPush
}
