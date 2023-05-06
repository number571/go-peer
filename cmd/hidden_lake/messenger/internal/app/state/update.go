package state

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func (p *sState) updateClientState(pStateValue *SStorageState) error {
	if err := p.updateClientPrivKey(pStateValue); err != nil {
		return err
	}

	if err := p.updateClientFriends(pStateValue); err != nil {
		return err
	}

	if err := p.updateClientConnections(pStateValue); err != nil {
		return err
	}

	return nil
}

func (p *sState) updateClientPrivKey(pStateValue *SStorageState) error {
	hlsClient := p.GetClient().Service()

	if err := p.clearClientPrivKey(); err != nil {
		return err
	}

	privKey := asymmetric.LoadRSAPrivKey(pStateValue.FPrivKey)
	if privKey == nil {
		return fmt.Errorf("private key is null")
	}

	return hlsClient.SetPrivKey(privKey)
}

func (p *sState) updateClientFriends(pStateValue *SStorageState) error {
	client := p.GetClient().Service()

	if err := p.clearClientFriends(); err != nil {
		return err
	}

	for aliasName, pubKeyString := range pStateValue.FFriends {
		pubKey := asymmetric.LoadRSAPubKey(pubKeyString)
		if err := client.AddFriend(aliasName, pubKey); err != nil {
			return err
		}
	}

	return nil
}

func (p *sState) updateClientConnections(pStateValue *SStorageState) error {
	client := p.GetClient().Service()

	if err := p.clearClientConnections(); err != nil {
		return err
	}

	for _, conn := range pStateValue.FConnections {
		if err := client.AddConnection(conn); err != nil {
			return err
		}
	}

	return nil
}

func (p *sState) updateClientTraffic(pStateValue *SStorageState) error {
	hlsClient := p.GetClient().Service()
	hltClient := p.GetClient().Traffic()

	hashes, err := hltClient.GetHashes()
	if err != nil {
		return err
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
