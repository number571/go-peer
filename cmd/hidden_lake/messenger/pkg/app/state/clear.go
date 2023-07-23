package state

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
)

func (p *sStateManager) clearClientState() error {
	if err := p.clearClientPrivKey(); err != nil {
		return errors.WrapError(err, "clear client private key")
	}

	if err := p.clearClientFriends(); err != nil {
		return errors.WrapError(err, "clear client friends")
	}

	if err := p.clearClientConnections(); err != nil {
		return errors.WrapError(err, "clear client connections")
	}

	return nil
}

func (p *sStateManager) clearClientPrivKey() error {
	client := p.GetClient().Service()

	pseudoPrivKey := asymmetric.NewRSAPrivKey(p.fConfig.GetKeySizeBits())
	if err := client.SetPrivKey(pseudoPrivKey); err != nil {
		return errors.WrapError(err, "set pseudo private key")
	}
	return nil
}

func (p *sStateManager) clearClientFriends() error {
	client := p.GetClient().Service()

	friends, err := client.GetFriends()
	if err != nil {
		return errors.WrapError(err, "get friends")
	}

	for aliasName := range friends {
		if err := client.DelFriend(aliasName); err != nil {
			return errors.WrapError(err, "del friend")
		}
	}

	return nil
}

func (p *sStateManager) clearClientConnections() error {
	client := p.GetClient().Service()

	connects, err := client.GetConnections()
	if err != nil {
		return errors.WrapError(err, "get connections")
	}
	for _, conn := range connects {
		if err := client.DelConnection(conn); err != nil {
			return errors.WrapError(err, "del connection")
		}
	}

	return nil
}
