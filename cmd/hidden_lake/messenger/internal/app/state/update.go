package state

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func (s *sState) updateClientState(stateValue *SStorageState) error {
	if err := s.updateClientPrivKey(stateValue); err != nil {
		return err
	}

	if err := s.updateClientFriends(stateValue); err != nil {
		return err
	}

	return s.updateClientConnections(stateValue)
}

func (s *sState) updateClientPrivKey(stateValue *SStorageState) error {
	client := s.GetClient()

	if err := s.clearClientPrivKey(); err != nil {
		return err
	}

	privKey := asymmetric.LoadRSAPrivKey(stateValue.FPrivKey)
	if privKey == nil {
		return fmt.Errorf("private key is null")
	}

	return client.SetPrivKey(privKey)
}

func (s *sState) updateClientFriends(stateValue *SStorageState) error {
	client := s.GetClient()

	if err := s.clearClientFriends(); err != nil {
		return err
	}

	for aliasName, pubKeyString := range stateValue.FFriends {
		pubKey := asymmetric.LoadRSAPubKey(pubKeyString)
		if err := client.AddFriend(aliasName, pubKey); err != nil {
			return err
		}
	}

	return nil
}

func (s *sState) updateClientConnections(stateValue *SStorageState) error {
	client := s.GetClient()

	if err := s.clearClientConnections(); err != nil {
		return err
	}

	for _, conn := range stateValue.FConnections {
		if err := client.AddConnection(conn); err != nil {
			return err
		}
	}

	return nil
}
