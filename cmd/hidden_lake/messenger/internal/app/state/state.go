package state

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/entropy"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

var (
	_ IState  = &sState{}
	_ iClient = &sClient{}
)

type sState struct {
	fMutex    sync.Mutex
	fHashLP   []byte
	fStorage  storage.IKVStorage
	fDatabase database.IWrapperDB
	fClient   *sClient
	fPathTo   string
}

type sClient struct {
	fService hls_client.IClient
	fTraffic hlt_client.IClient
}

func NewState(
	pStorage storage.IKVStorage,
	pDatabase database.IWrapperDB,
	pHlsClient hls_client.IClient,
	pHltClient hlt_client.IClient,
	pPathTo string,
) IState {
	return &sState{
		fStorage:  pStorage,
		fDatabase: pDatabase,
		fClient: &sClient{
			fService: pHlsClient,
			fTraffic: pHltClient,
		},
		fPathTo: pPathTo,
	}
}

func (p *sClient) Service() hls_client.IClient {
	return p.fService
}

func (p *sClient) Traffic() hlt_client.IClient {
	return p.fTraffic
}

func (p *sState) GetClient() iClient {
	return p.fClient
}

func (p *sState) GetKVStorage() storage.IKVStorage {
	return p.fStorage
}

func (p *sState) GetWrapperDB() database.IWrapperDB {
	return p.fDatabase
}

func (p *sState) GetTemplate() *STemplateState {
	return &STemplateState{
		FAuthorized: p.IsActive(),
	}
}

func (p *sState) CreateState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	if _, err := p.GetKVStorage().Get(pHashLP); err == nil {
		return errors.NewError("state already exists")
	}
	if err := p.newStorageState(pHashLP, pPrivKey); err != nil {
		return errors.WrapError(err, "new storage state")
	}
	return nil
}

func (p *sState) UpdateState(pHashLP []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.IsActive() {
		return errors.NewError("state already exists")
	}

	stateValue, err := p.getStorageState(pHashLP)
	if err != nil {
		return errors.WrapError(err, "get storage state")
	}

	entropyBooster := entropy.NewEntropyBooster(hlm_settings.CWorkForKeys, []byte{5, 7, 1})
	db, err := database.NewKeyValueDB(
		fmt.Sprintf("%s/%s", p.fPathTo, hlm_settings.CPathDB),
		entropyBooster.BoostEntropy(pHashLP),
	)
	if err != nil {
		return errors.WrapError(err, "open KV database")
	}

	p.GetWrapperDB().Set(db)
	if err := p.updateClientState(stateValue); err != nil {
		return errors.WrapError(err, "update client state")
	}
	p.fHashLP = pHashLP

	_ = p.updateClientTraffic(stateValue) // connect to HLT
	return nil
}

func (p *sState) ClearActiveState() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.IsActive() {
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

func (p *sState) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
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

func (p *sState) DelFriend(pAliasName string) error {
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

func (p *sState) AddConnection(pConnect string) error {
	err := p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = append(
				storageValue.FConnections,
				pConnect,
			)
		},
	)
	if err != nil {
		return errors.WrapError(err, "add connection (state updater)")
	}
	return nil
}

func (p *sState) DelConnection(pConnect string) error {
	err := p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = remove(
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

func (p *sState) IsActive() bool {
	return p.fHashLP != nil
}

func (p *sState) newStorageState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	stateValueBytes := encoding.Serialize(&SStorageState{
		FPrivKey: pPrivKey.ToString(),
	}, false)
	if err := p.GetKVStorage().Set(pHashLP, stateValueBytes); err != nil {
		return errors.WrapError(err, "new storage state")
	}
	return nil
}

func (p *sState) setStorageState(pStateValue *SStorageState) error {
	stateValueBytes := encoding.Serialize(pStateValue, false)
	if err := p.GetKVStorage().Set(p.fHashLP, stateValueBytes); err != nil {
		return errors.WrapError(err, "update storage state")
	}
	return nil
}

func (p *sState) getStorageState(pHashLP []byte) (*SStorageState, error) {
	stateValueBytes, err := p.GetKVStorage().Get(pHashLP)
	if err != nil {
		return nil, errors.WrapError(err, "get storage state")
	}

	var stateValue = new(SStorageState)
	if err := encoding.Deserialize(stateValueBytes, stateValue); err != nil {
		return nil, errors.WrapError(err, "deserialize state")
	}

	return stateValue, nil
}

func remove(pSlice []string, pElem string) []string {
	for i, sElem := range pSlice {
		if pElem == sElem {
			return append(pSlice[:i], pSlice[i+1:]...)
		}
	}
	return pSlice
}
