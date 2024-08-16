package client

import (
	"bytes"
	"errors"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fSettings   ISettings
	fStructSize uint64
}

// Create client by private key as identification.
// Handle function is used when the network exists.
func NewClient(pSett ISettings) IClient {
	client := &sClient{
		fSettings: pSett,
	}

	tempKey := random.NewCSPRNG().GetBytes(symmetric.CAESKeySize)
	encMsg, err := client.encryptWithParams(tempKey, []byte{}, 0)
	if err != nil {
		panic(err)
	}

	client.fStructSize = uint64(len(encMsg))
	if limit := client.GetMessageLimit(); limit == 0 {
		panic("the message size is lower than struct size")
	}

	return client
}

// Message is raw bytes of body payload without payload head.
func (p *sClient) GetMessageLimit() uint64 {
	maxMsgSize := p.fSettings.GetMessageSizeBytes()
	if maxMsgSize <= p.fStructSize {
		return 0
	}
	return maxMsgSize - p.fStructSize
}

// Get settings from client object.
func (p *sClient) GetSettings() ISettings {
	return p.fSettings
}

func (p *sClient) MessageIsValid(pMsg []byte) bool {
	return uint64(len(pMsg)) == p.fSettings.GetMessageSizeBytes()
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (p *sClient) EncryptMessage(pKey []byte, pMsg []byte) ([]byte, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pMsg))
	)
	if resultSize > msgLimitSize {
		return nil, ErrLimitMessageSize
	}
	return p.encryptWithParams(pKey, pMsg, msgLimitSize-resultSize)
}

func (p *sClient) encryptWithParams(pKey []byte, pMsg []byte, pPadd uint64) ([]byte, error) {
	data := joiner.NewBytesJoiner32([][]byte{pMsg, random.NewCSPRNG().GetBytes(pPadd)})
	hash := hashing.NewHMACSHA256Hasher(pKey, data).ToBytes()
	ciph := symmetric.NewAESCipher(pKey)
	return ciph.EncryptBytes(bytes.Join([][]byte{hash, data}, []byte{})), nil
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pKey []byte, pMsg []byte) ([]byte, error) {
	if !p.MessageIsValid(pMsg) {
		return nil, errors.New("message is invalid") // nolint: goerr113
	}

	ciph := symmetric.NewAESCipher(pKey)
	dmsg := ciph.DecryptBytes(pMsg)
	if dmsg == nil {
		return nil, errors.New("decrypt message") // nolint: goerr113
	}

	var (
		hash = dmsg[:hashing.CSHA256Size]
		data = dmsg[hashing.CSHA256Size:]
	)

	check := hashing.NewHMACSHA256Hasher(pKey, data).ToBytes()
	if !bytes.Equal(check, hash) {
		return nil, errors.New("hash message") // nolint: goerr113
	}

	dataWrapper, err := joiner.LoadBytesJoiner32(data)
	if err != nil || len(dataWrapper) != 2 {
		return nil, errors.New("payload") // nolint: goerr113
	}

	return dataWrapper[0], nil
}
