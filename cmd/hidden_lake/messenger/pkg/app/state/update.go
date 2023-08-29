package state

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

func (p *sStateManager) updateClientState(pStateValue *SStorageState) error {
	if err := p.updateClientPrivKey(pStateValue); err != nil {
		return errors.WrapError(err, "update client private key")
	}

	if err := p.updateClientFriends(pStateValue); err != nil {
		return errors.WrapError(err, "update client friends")
	}

	if err := p.updateClientNetworkKey(pStateValue); err != nil {
		return errors.WrapError(err, "update client connections")
	}

	return nil
}

func (p *sStateManager) updateClientPrivKey(pStateValue *SStorageState) error {
	hlsClient := p.GetClient()

	_, ephPubKey, err := hlsClient.GetPubKey()
	if err != nil {
		return errors.WrapError(err, "get public key from node (update)")
	}

	privKey := asymmetric.LoadRSAPrivKey(pStateValue.FPrivKey)
	if privKey == nil {
		return errors.NewError("private key is null (update)")
	}

	if err := hlsClient.SetPrivKey(privKey, ephPubKey); err != nil {
		return errors.WrapError(err, "set private key (update)")
	}
	return nil
}

func (p *sStateManager) updateClientFriends(pStateValue *SStorageState) error {
	client := p.GetClient()

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

func (p *sStateManager) updateClientNetworkKey(pStateValue *SStorageState) error {
	client := p.GetClient()

	if err := client.SetNetworkKey(pStateValue.FNetworkKey); err != nil {
		return errors.WrapError(err, "update client network key")
	}

	return nil
}

func (p *sStateManager) updateClientTraffic(pStateValue *SStorageState) {
	for _, conn := range p.fConfig.GetBackupConnections() {
		go p.handleMessages(conn)
	}
}

func (p *sStateManager) handleMessages(pConn string) {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", pConn),
			&http.Client{Timeout: time.Minute},
			message.NewSettings(&message.SSettings{
				FWorkSizeBits:     p.fConfig.GetWorkSizeBits(),
				FMessageSizeBytes: p.fConfig.GetMessageSizeBytes(),
			}),
		),
	)

	hashes, err := hltClient.GetHashes()
	if err != nil {
		return
	}

	hlsClient := p.GetClient()
	for i, hash := range hashes {
		if uint64(i) >= p.fConfig.GetMessagesCapacity() {
			break
		}
		msg, err := hltClient.GetMessage(hash)
		if err != nil {
			continue
		}
		bytesHash := encoding.HexDecode(hash)
		if !bytes.Equal(msg.GetBody().GetHash(), bytesHash) {
			break
		}
		if err := hlsClient.HandleMessage(msg); err != nil {
			continue
		}
	}
}
