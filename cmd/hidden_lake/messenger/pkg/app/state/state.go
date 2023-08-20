package state

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/stringtools"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
)

var (
	_ IStateManager = &sStateManager{}
)

type sStateManager struct {
	fMutex    sync.Mutex
	fHashLP   []byte
	fPrivKey  asymmetric.IPrivKey
	fConfig   config.IConfig
	fPathTo   string
	fStorage  storage.IKVStorage
	fDatabase database.IWrapperDB
	fClient   hls_client.IClient
}

func NewStateManager(
	pConfig config.IConfig,
	pPathTo string,
) IStateManager {
	stg, err := initCryptoStorage(pConfig, pPathTo)
	if err != nil {
		panic(err)
	}
	return &sStateManager{
		fConfig:   pConfig,
		fPathTo:   pPathTo,
		fStorage:  stg,
		fDatabase: database.NewWrapperDB(),
		fClient: hls_client.NewClient(
			hls_client.NewBuilder(),
			hls_client.NewRequester(
				fmt.Sprintf("http://%s", pConfig.GetConnection()),
				&http.Client{Timeout: time.Minute},
			),
		),
	}
}

func (p *sStateManager) GetConfig() config.IConfig {
	return p.fConfig
}

func (p *sStateManager) GetClient() hls_client.IClient {
	return p.fClient
}

func (p *sStateManager) GetWrapperDB() database.IWrapperDB {
	return p.fDatabase
}

func (p *sStateManager) GetTemplate() *STemplateState {
	return &STemplateState{
		FLanguage:   p.fConfig.GetLanguage(),
		FAuthorized: p.StateIsActive(),
	}
}

func (p *sStateManager) GetPrivKey() asymmetric.IPrivKey {
	if !p.StateIsActive() {
		return nil
	}
	return p.fPrivKey
}

func (p *sStateManager) IsMyPubKey(pPubKey asymmetric.IPubKey) bool {
	if !p.StateIsActive() {
		return false
	}
	myPubKey := p.fPrivKey.GetPubKey()
	if myPubKey == nil || pPubKey == nil {
		return false
	}
	return myPubKey.ToString() == pPubKey.ToString()
}

func (p *sStateManager) CreateState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	if _, err := p.fStorage.Get(pHashLP); err == nil {
		return errors.NewError("state already exists")
	}
	if err := p.newStorageState(pHashLP, pPrivKey); err != nil {
		return errors.WrapError(err, "new storage state")
	}
	return nil
}

func (p *sStateManager) OpenState(pHashLP []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.StateIsActive() {
		return errors.NewError("state already exists")
	}

	stateValue, err := p.getStorageState(pHashLP)
	if err != nil {
		return errors.WrapError(err, "get storage state")
	}

	db, err := database.NewKeyValueDB(
		fmt.Sprintf("%s/%s", p.fPathTo, hlm_settings.CPathDB),
		encoding.HexEncode(pHashLP),
	)
	if err != nil {
		return errors.WrapError(err, "open KV database")
	}

	p.GetWrapperDB().Set(db)
	if err := p.updateClientState(stateValue); err != nil {
		return errors.WrapError(err, "update client state")
	}

	privKey := asymmetric.LoadRSAPrivKey(stateValue.FPrivKey)
	if privKey == nil {
		return errors.NewError("private key is null (open state)")
	}

	p.fHashLP = pHashLP
	p.fPrivKey = privKey

	p.updateClientTraffic(stateValue)
	return nil
}

func (p *sStateManager) CloseState() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.StateIsActive() {
		return errors.NewError("state does not exist")
	}

	p.fHashLP = nil

	if err := p.GetWrapperDB().Close(); err != nil {
		return errors.WrapError(err, "close database")
	}

	if err := p.clearClientState(); err != nil {
		return errors.WrapError(err, "clear client state")
	}

	return nil
}

func (p *sStateManager) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
	err := p.stateUpdater(
		p.updateClientFriends,
		func(storageValue *SStorageState) {
			storageValue.FFriends[pAliasName] = pPubKey.ToString()
		},
	)
	if err != nil {
		return errors.WrapError(err, "add friend (state updater)")
	}
	return nil
}

func (p *sStateManager) DelFriend(pAliasName string) error {
	err := p.stateUpdater(
		p.updateClientFriends,
		func(storageValue *SStorageState) {
			delete(storageValue.FFriends, pAliasName)
		},
	)
	if err != nil {
		return errors.WrapError(err, "del friend (state updater)")
	}
	return nil
}

func (p *sStateManager) SetNetworkKey(pNetworkKey string) error {
	err := p.stateUpdater(
		p.updateClientNetworkKey,
		func(storageValue *SStorageState) {
			storageValue.FNetworkKey = pNetworkKey
		},
	)
	if err != nil {
		return errors.WrapError(err, "set network key (state updater)")
	}
	return nil
}

func (p *sStateManager) AddConnection(pAddress string) error {
	err := p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = stringtools.UniqAppendToSlice(
				storageValue.FConnections,
				pAddress,
			)
		},
	)
	if err != nil {
		return errors.WrapError(err, "add connection (state updater)")
	}
	return nil
}

func (p *sStateManager) DelConnection(pConnect string) error {
	err := p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = stringtools.DeleteFromSlice(
				storageValue.FConnections,
				pConnect,
			)
		},
	)
	if err != nil {
		return errors.WrapError(err, "del connection (state updater)")
	}
	return nil
}

func (p *sStateManager) StateIsActive() bool {
	return p.fHashLP != nil
}

func (p *sStateManager) newStorageState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	stateValueBytes := encoding.Serialize(
		&SStorageState{
			FPrivKey: pPrivKey.ToString(),
		},
		false,
	)
	if err := p.fStorage.Set(pHashLP, stateValueBytes); err != nil {
		return errors.WrapError(err, "new storage state")
	}
	return nil
}

func (p *sStateManager) setStorageState(pStateValue *SStorageState) error {
	stateValueBytes := encoding.Serialize(pStateValue, false)
	if err := p.fStorage.Set(p.fHashLP, stateValueBytes); err != nil {
		return errors.WrapError(err, "update storage state")
	}
	return nil
}

func (p *sStateManager) getStorageState(pHashLP []byte) (*SStorageState, error) {
	stateValueBytes, err := p.fStorage.Get(pHashLP)
	if err != nil {
		return nil, errors.WrapError(err, "get storage state bytes")
	}

	var stateValue = new(SStorageState)
	if err := encoding.Deserialize(stateValueBytes, stateValue); err != nil {
		return nil, errors.WrapError(err, "deserialize state")
	}

	return stateValue, nil
}
