package state

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/utils"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type STemplateState struct {
	FLanguage   utils.ILanguage
	FAuthorized bool
}

type SStorageState struct {
	FNetworkKey  string            `json:"network_key"`
	FPrivKey     string            `json:"priv_key"`
	FConnections []string          `json:"connections"`
	FFriends     map[string]string `json:"friends"`
}

type IStateManager interface {
	GetConfig() config.IConfig

	IsMyPubKey(asymmetric.IPubKey) bool
	GetPrivKey() asymmetric.IPrivKey

	CreateState([]byte, asymmetric.IPrivKey) error

	StateIsActive() bool
	OpenState([]byte) error
	CloseState() error

	GetClient() hls_client.IClient
	GetWrapperDB() database.IWrapperDB
	GetTemplate() *STemplateState

	SetNetworkKey(string) error

	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	AddConnection(string) error
	DelConnection(string) error
}
