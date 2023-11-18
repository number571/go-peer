package client

import (
	"bytes"
	"fmt"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fSettings   message.ISettings
	fPrivKey    asymmetric.IPrivKey
	fStructSize uint64
}

// Create client by private key as identification.
// Handle function is used when the network exists.
func NewClient(pSett message.ISettings, pPrivKey asymmetric.IPrivKey) IClient {
	client := &sClient{
		fSettings: pSett,
		fPrivKey:  pPrivKey,
	}

	encMsg := client.encryptWithParams(
		client.GetPubKey(),
		payload.NewPayload(0, []byte{}),
		0,
	)

	client.fStructSize = uint64(len(encMsg.ToBytes()))
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

// Get public key from client object.
func (p *sClient) GetPubKey() asymmetric.IPubKey {
	return p.fPrivKey.GetPubKey()
}

// Get private key from client object.
func (p *sClient) GetPrivKey() asymmetric.IPrivKey {
	return p.fPrivKey
}

// Get settings from client object.
func (p *sClient) GetSettings() message.ISettings {
	return p.fSettings
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (p *sClient) EncryptPayload(pRecv asymmetric.IPubKey, pPld payload.IPayload) (message.IMessage, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pPld.GetBody()))
	)

	if resultSize > msgLimitSize {
		return nil, errors.NewError(fmt.Sprintf(
			"limit of message size without hex encoding = %d bytes < current payload size with additional padding = %d bytes",
			msgLimitSize,
			resultSize,
		))
	}

	return p.encryptWithParams(
		pRecv,
		pPld,
		msgLimitSize-resultSize,
	), nil
}

func (p *sClient) encryptWithParams(pRecv asymmetric.IPubKey, pPld payload.IPayload, pPadd uint64) message.IMessage {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.GetBytes(symmetric.CAESKeySize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	payloadBytes := pPld.ToBytes()
	doublePayload := payload.NewPayload(
		uint64(len(payloadBytes)),
		bytes.Join(
			[][]byte{
				payloadBytes,
				rand.GetBytes(pPadd),
			},
			[]byte{},
		),
	)

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			p.GetPubKey().GetAddress().ToBytes(),
			pRecv.GetAddress().ToBytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).ToBytes()

	cipher := symmetric.NewAESCipher(session)
	return &message.SMessage{
		FPubKey:  encoding.HexEncode(cipher.EncryptBytes(p.GetPubKey().ToBytes())),
		FEncKey:  encoding.HexEncode(pRecv.EncryptBytes(session)),
		FSalt:    encoding.HexEncode(cipher.EncryptBytes(salt)),
		FHash:    encoding.HexEncode(hash),
		FSign:    encoding.HexEncode(cipher.EncryptBytes(p.fPrivKey.SignBytes(hash))),
		FPayload: cipher.EncryptBytes(doublePayload.ToBytes()), // JSON field to raw Body (no need HEX encode)
	}
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pMsg message.IMessage) (asymmetric.IPubKey, payload.IPayload, error) {
	// Initial check.
	if pMsg == nil || !pMsg.IsValid(p.fSettings) {
		return nil, nil, errors.NewError("got invalid message")
	}

	// Decrypt session key by private key of receiver.
	session := p.fPrivKey.DecryptBytes(pMsg.GetEncKey())
	if session == nil {
		return nil, nil, errors.NewError("failed decrypt session key")
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := symmetric.NewAESCipher(session)
	publicBytes := cipher.DecryptBytes(pMsg.GetPubKey())
	if publicBytes == nil {
		return nil, nil, errors.NewError("failed decrypt public key")
	}

	// Load public key and check standart size.
	pubKey := asymmetric.LoadRSAPubKey(publicBytes)
	if pubKey == nil {
		return nil, nil, errors.NewError("failed load public key")
	}
	if pubKey.GetSize() != p.GetPubKey().GetSize() {
		return nil, nil, errors.NewError("invalid public key size")
	}

	// Decrypt main data of message by session key.
	firstPayload := pMsg.GetPayload()
	if firstPayload == nil {
		return nil, nil, errors.NewError("failed decode payload")
	}
	doublePayloadBytes := cipher.DecryptBytes(firstPayload.ToBytes())
	if doublePayloadBytes == nil {
		return nil, nil, errors.NewError("failed decrypt double payload")
	}
	doublePayload := payload.LoadPayload(doublePayloadBytes)
	if doublePayload == nil {
		return nil, nil, errors.NewError("failed load double payload")
	}

	// Decrypt salt.
	salt := cipher.DecryptBytes(pMsg.GetSalt())
	if salt == nil {
		return nil, nil, errors.NewError("failed decrypt salt")
	}

	// Check received hash and generated hash.
	check := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			pubKey.GetAddress().ToBytes(),
			p.GetPubKey().GetAddress().ToBytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, pMsg.GetHash()) {
		return nil, nil, errors.NewError("invalid msg hash")
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.DecryptBytes(pMsg.GetSign())
	if sign == nil {
		return nil, nil, errors.NewError("failed decrypt sign")
	}
	if !pubKey.VerifyBytes(pMsg.GetHash(), sign) {
		return nil, nil, errors.NewError("invalid msg sign")
	}

	// Remove random bytes and get main data
	mustLen := doublePayload.GetHead()
	if mustLen > uint64(len(doublePayload.GetBody())) {
		return nil, nil, errors.NewError("invalid size of payload")
	}
	pld := payload.LoadPayload(doublePayload.GetBody()[:mustLen])
	if pld == nil {
		return nil, nil, errors.NewError("invalid load payload")
	}

	// Return decrypted message with title
	return pubKey, pld, nil
}
