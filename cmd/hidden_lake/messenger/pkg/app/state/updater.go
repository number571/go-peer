package state

import "github.com/number571/go-peer/pkg/errors"

func (p *sStateManager) stateUpdater(
	clientUpdater func(storageValue *SStorageState) error,
	middleWare func(storageValue *SStorageState),
) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.StateIsActive() {
		return errors.NewError("state does not exist")
	}

	oldStorageValue, err := p.getStorageState(p.fHashLP)
	if err != nil {
		return errors.WrapError(err, "get old storage state")
	}

	newStorageValue := copyStorageState(oldStorageValue)
	middleWare(newStorageValue)

	if err := clientUpdater(newStorageValue); err == nil {
		if err := p.setStorageState(newStorageValue); err != nil {
			return errors.WrapError(err, "set new storage state")
		}
		return nil
	}

	if err := p.setStorageState(oldStorageValue); err != nil {
		return errors.WrapError(err, "update state (old -> new)")
	}
	return nil
}

func copyStorageState(pStorageValue *SStorageState) *SStorageState {
	copyStorageValue := &SStorageState{
		FPrivKey: pStorageValue.FPrivKey,
		FFriends: make(map[string]string, len(pStorageValue.FFriends)),
	}

	for aliasName, pubKey := range pStorageValue.FFriends {
		copyStorageValue.FFriends[aliasName] = pubKey
	}

	return copyStorageValue
}
