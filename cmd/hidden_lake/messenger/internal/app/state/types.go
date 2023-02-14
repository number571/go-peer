package state

import (
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage"
)

type STemplateState struct {
	FAuthorized bool
}

type SStorageState struct {
	FPrivKey     string            `json:"priv_key"`
	FConnections []string          `json:"connections"`
	FFriends     map[string]string `json:"friends"`
}

type IState interface {
	GetClient() iClient
	GetStorage() storage.IKeyValueStorage
	GetWrapperDB() database.IWrapperDB

	IsActive() bool
	GetTemplate() *STemplateState

	AddFriend(string, asymmetric.IPubKey) error
	DelFriend(string) error

	AddConnection(string) error
	DelConnection(string) error

	CreateState([]byte, asymmetric.IPrivKey) error
	UpdateState([]byte) error
	ClearActiveState() error
}

type iClient interface {
	Service() hls_client.IClient
	Traffic() hlt_client.IClient
}
