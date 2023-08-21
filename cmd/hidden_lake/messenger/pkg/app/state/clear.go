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
	hlsClient := p.GetClient()

	_, ephPubKey, err := hlsClient.GetPubKey()
	if err != nil {
		return errors.WrapError(err, "get public key from node (clear)")
	}

	pseudoPrivKey := asymmetric.NewRSAPrivKey(p.fConfig.GetKeySizeBits())
	if err := hlsClient.SetPrivKey(pseudoPrivKey, ephPubKey); err != nil {
		return errors.WrapError(err, "set pseudo private key (clear)")
	}
	return nil
}

func (p *sStateManager) clearClientFriends() error {
	client := p.GetClient()

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
	client := p.GetClient()

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
