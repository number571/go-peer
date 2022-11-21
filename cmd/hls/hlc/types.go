package hlc

import (
	"github.com/number571/go-peer/modules/crypto/asymmetric"

	hls_network "github.com/number571/go-peer/cmd/hls/network"
	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

type IClient interface {
	PubKey() (asymmetric.IPubKey, error)

	GetOnlines() ([]string, error)
	DelOnline(string) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	GetConnections() ([]string, error)
	AddConnection(string) error
	DelConnection(string) error

	Broadcast(asymmetric.IPubKey, hls_network.IRequest) error
	Request(asymmetric.IPubKey, hls_network.IRequest) ([]byte, error)
}

type IRequester interface {
	PubKey() (asymmetric.IPubKey, error)

	GetOnlines() ([]string, error)
	DelOnline(*hls_settings.SConnect) error

	GetFriends() (map[string]asymmetric.IPubKey, error)
	AddFriend(*hls_settings.SFriend) error
	DelFriend(*hls_settings.SFriend) error

	GetConnections() ([]string, error)
	AddConnection(*hls_settings.SConnect) error
	DelConnection(*hls_settings.SConnect) error

	Broadcast(*hls_settings.SPush) error
	Request(*hls_settings.SPush) ([]byte, error)
}

type IBuilder interface {
	Connect(string) *hls_settings.SConnect
	Friend(string, asymmetric.IPubKey) *hls_settings.SFriend
	Push(asymmetric.IPubKey, hls_network.IRequest) *hls_settings.SPush
}
