package state

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
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

	_ = p.updateClientTraffic(pStateValue)
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

func (p *sState) updateClientTraffic(pStateValue *SStorageState) error {
	hltClient := p.GetClient().Traffic()

	privKey := asymmetric.LoadRSAPrivKey(pStateValue.FPrivKey)
	if privKey == nil {
		return fmt.Errorf("private key is null")
	}

	hashes, err := hltClient.GetHashes()
	if err != nil {
		return err
	}

	for _, hash := range hashes {
		msg, err := hltClient.GetMessage(hash)
		if err != nil {
			continue
		}

		c := hls_settings.InitClient(privKey)
		pubKey, pld, err := c.DecryptMessage(msg)
		if err != nil {
			continue
		}

		inFriends := false
		for _, pubKeyString := range pStateValue.FFriends {
			fPubKey := asymmetric.LoadRSAPubKey(pubKeyString)
			if pubKey.GetAddress().ToString() == fPubKey.GetAddress().ToString() {
				inFriends = true
				break
			}
		}

		if !inFriends {
			continue
		}

		loadReq, err := request.LoadRequest(pld.GetBody())
		if err != nil {
			continue
		}

		strMsg := strings.TrimSpace(string(loadReq.Body()))
		if len(strMsg) == 0 {
			continue
		}

		rel := database.NewRelation(privKey.GetPubKey(), pubKey)
		dbMsg := database.NewMessage(true, strMsg, msg.GetBody().GetHash())

		db := p.GetWrapperDB().Get()
		if err := db.Push(rel, dbMsg); err != nil {
			continue
		}
	}

	return nil
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
