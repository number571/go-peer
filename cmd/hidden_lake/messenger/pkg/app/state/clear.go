package state

import (
	"github.com/number571/go-peer/pkg/errors"
)

func (p *sStateManager) clearClientState() error {
	if err := p.clearClientPrivKey(); err != nil {
		return errors.WrapError(err, "clear client private key")
	}

	if err := p.clearClientFriends(); err != nil {
		return errors.WrapError(err, "clear client friends")
	}

	return nil
}

func (p *sStateManager) clearClientPrivKey() error {
	hlsClient := p.GetClient()

	if err := hlsClient.ResetPrivKey(); err != nil {
		return errors.WrapError(err, "reset private key (clear)")
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
