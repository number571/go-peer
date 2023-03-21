package state

import (
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func (p *sState) clearClientState() error {
	if err := p.clearClientPrivKey(); err != nil {
		return err
	}

	if err := p.clearClientFriends(); err != nil {
		return err
	}

	return p.clearClientConnections()
}

func (p *sState) clearClientPrivKey() error {
	client := p.GetClient().Service()

	pseudoPrivKey := asymmetric.NewRSAPrivKey(pkg_settings.CAKeySize)
	return client.SetPrivKey(pseudoPrivKey)
}

func (p *sState) clearClientFriends() error {
	client := p.GetClient().Service()

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

func (p *sState) clearClientConnections() error {
	client := p.GetClient().Service()

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
