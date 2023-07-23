package state

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/errors"
)

func (p *sStateManager) updateClientState(pStateValue *SStorageState) error {
	if err := p.updateClientPrivKey(pStateValue); err != nil {
		return errors.WrapError(err, "update client private key")
	}

	if err := p.updateClientFriends(pStateValue); err != nil {
		return errors.WrapError(err, "update client friends")
	}

	if err := p.updateClientConnections(pStateValue); err != nil {
		return errors.WrapError(err, "update client connections")
	}

	return nil
}

func (p *sStateManager) updateClientPrivKey(pStateValue *SStorageState) error {
	hlsClient := p.GetClient().Service()

	if err := p.clearClientPrivKey(); err != nil {
		return errors.WrapError(err, "clear client private key")
	}

	privKey := asymmetric.LoadRSAPrivKey(pStateValue.FPrivKey)
	if privKey == nil {
		return errors.NewError("private key is null")
	}

	if err := hlsClient.SetPrivKey(privKey); err != nil {
		return errors.WrapError(err, "set private key")
	}
	return nil
}

func (p *sStateManager) updateClientFriends(pStateValue *SStorageState) error {
	client := p.GetClient().Service()

	if err := p.clearClientFriends(); err != nil {
		return errors.WrapError(err, "clear client friends")
	}

	for aliasName, pubKeyString := range pStateValue.FFriends {
		pubKey := asymmetric.LoadRSAPubKey(pubKeyString)
		if err := client.AddFriend(aliasName, pubKey); err != nil {
			return errors.WrapError(err, "add friend")
		}
	}

	return nil
}

func (p *sStateManager) updateClientConnections(pStateValue *SStorageState) error {
	client := p.GetClient().Service()

	if err := p.clearClientConnections(); err != nil {
		return errors.WrapError(err, "clear client connections")
	}

	for _, conn := range pStateValue.FConnections {
		if err := client.AddConnection(conn); err != nil {
			return errors.WrapError(err, "add connections")
		}
	}

	return nil
}

func (p *sStateManager) updateClientTraffic(pStateValue *SStorageState) error {
	hlsClient := p.GetClient().Service()
	hltClient := p.GetClient().Traffic()

	hashes, err := hltClient.GetHashes()
	if err != nil {
		return errors.WrapError(err, "get hashes")
	}

	for _, hash := range hashes {
		msg, err := hltClient.GetMessage(hash)
		if err != nil {
			continue
		}
		if err := hlsClient.HandleMessage(msg); err != nil {
			continue
		}
	}

	return nil
}
