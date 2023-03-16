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
	storage storage.IKeyValueStorage,
	database database.IWrapperDB,
	hlsClient hls_client.IClient,
	hltClient hlt_client.IClient,
) IState {
	return &sState{
		fStorage:  storage,
		fDatabase: database,
		fClient: &sClient{
			fService: hlsClient,
			fTraffic: hltClient,
		},
	}
}

func (c *sClient) Service() hls_client.IClient {
	return c.fService
}

func (c *sClient) Traffic() hlt_client.IClient {
	return c.fTraffic
}

func (s *sState) GetClient() iClient {
	return s.fClient
}

func (s *sState) GetStorage() storage.IKeyValueStorage {
	return s.fStorage
}

func (s *sState) GetWrapperDB() database.IWrapperDB {
	return s.fDatabase
}

func (s *sState) GetTemplate() *STemplateState {
	return &STemplateState{
		FAuthorized: s.IsActive(),
	}
}

func (s *sState) CreateState(hashLP []byte, privKey asymmetric.IPrivKey) error {
	if _, err := s.GetStorage().Get(hashLP); err == nil {
		return fmt.Errorf("state already exists")
	}
	return s.newStorageState(hashLP, privKey)
}

func (s *sState) UpdateState(hashLP []byte) error {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	if s.IsActive() {
		return fmt.Errorf("state already exists")
	}

	stateValue, err := s.getStorageState(hashLP)
	if err != nil {
		return err
	}

	entropyBooster := entropy.NewEntropyBooster(hlm_settings.CWorkForKeys, []byte{5, 7, 1})
	s.GetWrapperDB().Set(database.NewKeyValueDB(
		hlm_settings.CPathDB,
		entropyBooster.BoostEntropy(hashLP),
	))

	if err := s.updateClientState(stateValue); err != nil {
		return err
	}

	s.fHashLP = hashLP
	return nil
}

func (s *sState) ClearActiveState() error {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	if !s.IsActive() {
		return fmt.Errorf("state does not exist")
	}

	if err := s.GetWrapperDB().Close(); err != nil {
		return err
	}

	if err := s.clearClientState(); err != nil {
		return err
	}

	s.fHashLP = nil
	return nil
}

func (s *sState) AddFriend(aliasName string, pubKey asymmetric.IPubKey) error {
	return s.stateUpdater(
		s.updateClientFriends,
		func(storageValue *SStorageState) {
			storageValue.FFriends[aliasName] = pubKey.ToString()
		},
	)
}

func (s *sState) DelFriend(aliasName string) error {
	return s.stateUpdater(
		s.updateClientFriends,
		func(storageValue *SStorageState) {
			delete(storageValue.FFriends, aliasName)
		},
	)
}

func (s *sState) AddConnection(connect string) error {
	return s.stateUpdater(
		s.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = append(
				storageValue.FConnections,
				connect,
			)
		},
	)
}

func (s *sState) DelConnection(connect string) error {
	return s.stateUpdater(
		s.updateClientConnections,
		func(storageValue *SStorageState) {
			storageValue.FConnections = remove(
				storageValue.FConnections,
				connect,
			)
		},
	)
}

func (s *sState) IsActive() bool {
	return s.fHashLP != nil
}

func (s *sState) newStorageState(hashLP []byte, privKey asymmetric.IPrivKey) error {
	stateValueBytes := encoding.Serialize(&SStorageState{
		FPrivKey: privKey.ToString(),
	})
	return s.GetStorage().Set(hashLP, stateValueBytes)
}

func (s *sState) setStorageState(stateValue *SStorageState) error {
	stateValueBytes := encoding.Serialize(stateValue)
	return s.GetStorage().Set(s.fHashLP, stateValueBytes)
}

func (s *sState) getStorageState(hashLP []byte) (*SStorageState, error) {
	stateValueBytes, err := s.GetStorage().Get(hashLP)
	if err != nil {
		return nil, err
	}

	var stateValue = new(SStorageState)
	if err := encoding.Deserialize(stateValueBytes, stateValue); err != nil {
		return nil, err
	}

	return stateValue, err
}

func remove(slice []string, elem string) []string {
	for i, sElem := range slice {
		if elem == sElem {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
