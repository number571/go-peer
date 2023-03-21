package state

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/entropy"
	"github.com/number571/go-peer/pkg/encoding"
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
	fStorage  storage.IKeyValueStorage
	fDatabase database.IWrapperDB
	fClient   *sClient
}

type sClient struct {
	fService hls_client.IClient
	fTraffic hlt_client.IClient
}

func NewState(
	pStorage storage.IKeyValueStorage,
	pDatabase database.IWrapperDB,
	pHlsClient hls_client.IClient,
	pHltClient hlt_client.IClient,
) IState {
	return &sState{
		fStorage:  pStorage,
		fDatabase: pDatabase,
		fClient: &sClient{
			fService: pHlsClient,
			fTraffic: pHltClient,
		},
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

func (p *sState) GetStorage() storage.IKeyValueStorage {
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
	if _, err := p.GetStorage().Get(pHashLP); err == nil {
		return fmt.Errorf("state already exists")
	}
	return p.newStorageState(pHashLP, pPrivKey)
}

func (p *sState) UpdateState(pHashLP []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.IsActive() {
		return fmt.Errorf("state already exists")
	}

	stateValue, err := p.getStorageState(pHashLP)
	if err != nil {
		return err
	}

	entropyBooster := entropy.NewEntropyBooster(hlm_settings.CWorkForKeys, []byte{5, 7, 1})
	p.GetWrapperDB().Set(database.NewKeyValueDB(
		hlm_settings.CPathDB,
		entropyBooster.BoostEntropy(pHashLP),
	))

	if err := p.updateClientState(stateValue); err != nil {
		return err
	}

	p.fHashLP = pHashLP
	return nil
}

func (p *sState) ClearActiveState() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.IsActive() {
		return fmt.Errorf("state does not exist")
	}

	p.fHashLP = nil

	if err := p.GetWrapperDB().Close(); err != nil {
		return err
	}

	if err := p.clearClientState(); err != nil {
		return err
	}

	return nil
}

func (p *sState) AddFriend(pAliasName string, pPubKey asymmetric.IPubKey) error {
	return p.stateUpdater(
		p.updateClientFriends,
		func(storageValue *SStorageState) {
			storageValue.FFriends[pAliasName] = pPubKey.ToString()
		},
	)
}

func (p *sState) DelFriend(pAliasName string) error {
	return p.stateUpdater(
		p.updateClientFriends,
		func(storageValue *SStorageState) {
			delete(storageValue.FFriends, pAliasName)
		},
	)
}

func (p *sState) AddConnection(pConnect string) error {
	return p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = append(
				storageValue.FConnections,
				pConnect,
			)
		},
	)
}

func (p *sState) DelConnection(pConnect string) error {
	return p.stateUpdater(
		p.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = remove(
				storageValue.FConnections,
				pConnect,
			)
		},
	)
}

func (p *sState) IsActive() bool {
	return p.fHashLP != nil
}

func (p *sState) newStorageState(pHashLP []byte, pPrivKey asymmetric.IPrivKey) error {
	stateValueBytes := encoding.Serialize(&SStorageState{
		FPrivKey: pPrivKey.ToString(),
	})
	return p.GetStorage().Set(pHashLP, stateValueBytes)
}

func (p *sState) setStorageState(pStateValue *SStorageState) error {
	stateValueBytes := encoding.Serialize(pStateValue)
	return p.GetStorage().Set(p.fHashLP, stateValueBytes)
}

func (p *sState) getStorageState(pHashLP []byte) (*SStorageState, error) {
	stateValueBytes, err := p.GetStorage().Get(pHashLP)
	if err != nil {
		return nil, err
	}

	var stateValue = new(SStorageState)
	if err := encoding.Deserialize(stateValueBytes, stateValue); err != nil {
		return nil, err
	}

	return stateValue, err
}

func remove(pSlice []string, pElem string) []string {
	for i, sElem := range pSlice {
		if pElem == sElem {
			return append(pSlice[:i], pSlice[i+1:]...)
		}
	}
	return pSlice
}
