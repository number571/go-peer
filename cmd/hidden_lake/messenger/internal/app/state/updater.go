package state

import "fmt"

func (s *sState) stateUpdater(
	clientUpdater func(storageValue *SStorageState) error,
	middleWare func(storageValue *SStorageState),
) error {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	if !s.IsActive() {
		return fmt.Errorf("state does not exist")
	}

	oldStorageValue, err := s.getStorageState(s.fHashLP)
	if err != nil {
		return err
	}

	newStorageValue := copyStorageState(oldStorageValue)
	middleWare(newStorageValue)

	if err := clientUpdater(newStorageValue); err == nil {
		return s.setStorageState(newStorageValue)
	}

	return s.setStorageState(oldStorageValue)
}

func copyStorageState(storageValue *SStorageState) *SStorageState {
	copyStorageValue := &SStorageState{
		FPrivKey:     storageValue.FPrivKey,
		FFriends:     make(map[string]string, len(storageValue.FFriends)),
		FConnections: make([]string, 0, len(storageValue.FConnections)),
	}

	for aliasName, pubKey := range storageValue.FFriends {
		copyStorageValue.FFriends[aliasName] = pubKey
	}

	copyStorageValue.FConnections = append(
		copyStorageValue.FConnections,
		storageValue.FConnections...,
	)

	return copyStorageValue
}
