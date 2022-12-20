package state

import (
	pkg_settings "github.com/number571/go-peer/cmd/hls/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func (s *sState) clearClientState() error {
	if err := s.clearClientPrivKey(); err != nil {
		return err
	}

	if err := s.clearClientFriends(); err != nil {
		return err
	}

	return s.clearClientConnections()
}

func (s *sState) clearClientPrivKey() error {
	client := s.GetClient()

	pseudoPrivKey := asymmetric.NewRSAPrivKey(pkg_settings.CAKeySize)
	return client.PrivKey(pseudoPrivKey)
}

func (s *sState) clearClientFriends() error {
	client := s.GetClient()

	friends, err := client.GetFriends()
	if err != nil {
		return err
	}

	for aliasName := range friends {
		if err := client.DelFriend(aliasName); err != nil {
			return err
		}
	}

	return nil
}

func (s *sState) clearClientConnections() error {
	client := s.GetClient()

	connects, err := client.GetConnections()
	if err != nil {
		return err
	}
	for _, conn := range connects {
		if err := client.DelConnection(conn); err != nil {
			return err
		}
	}

	return nil
}
