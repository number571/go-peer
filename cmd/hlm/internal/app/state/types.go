package state

import (
	"github.com/number571/go-peer/cmd/hlm/internal/database"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
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
	GetClient() hls_client.IClient
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
