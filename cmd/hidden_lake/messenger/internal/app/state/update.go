package state

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/database"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func (s *sState) updateClientState(stateValue *SStorageState) error {
	if err := s.updateClientPrivKey(stateValue); err != nil {
		return err
	}

	if err := s.updateClientFriends(stateValue); err != nil {
		return err
	}

	if err := s.updateClientConnections(stateValue); err != nil {
		return err
	}

	_ = s.updateClientTraffic(stateValue)
	return nil
}

func (s *sState) updateClientPrivKey(stateValue *SStorageState) error {
	hlsClient := s.GetClient().Service()

	if err := s.clearClientPrivKey(); err != nil {
		return err
	}

	privKey := asymmetric.LoadRSAPrivKey(stateValue.FPrivKey)
	if privKey == nil {
		return fmt.Errorf("private key is null")
	}

	return hlsClient.SetPrivKey(privKey)
}

func (s *sState) updateClientTraffic(stateValue *SStorageState) error {
	hltClient := s.GetClient().Traffic()

	privKey := asymmetric.LoadRSAPrivKey(stateValue.FPrivKey)
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
		for _, pubKeyString := range stateValue.FFriends {
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

		db := s.GetWrapperDB().Get()
		if err := db.Push(rel, dbMsg); err != nil {
			continue
		}
	}

	return nil
}

func (s *sState) updateClientFriends(stateValue *SStorageState) error {
	client := s.GetClient().Service()

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
	client := s.GetClient().Service()

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
