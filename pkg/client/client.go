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
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/utils"
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
	if pSett.GetKeySizeBits() != pPrivKey.GetSize() {
		panic("settings key size != got key size")
	}

	client := &sClient{
		fSettings: pSett,
		fPrivKey:  pPrivKey,
	}

	encMsg, err := client.encryptWithParams(
		client.GetPubKey(),
		payload.NewPayload(0, []byte{}),
		0,
	)
	if err != nil {
		panic(err)
	}

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
		return nil, utils.MergeErrors(
			ErrLimitMessageSize,
			fmt.Errorf(
				"limit of message size without hex encoding = %d bytes < current payload size with additional padding = %d bytes",
				msgLimitSize,
				resultSize,
			),
		)
	}

	return p.encryptWithParams(
		pRecv,
		pPld,
		msgLimitSize-resultSize,
	)
}

func (p *sClient) encryptWithParams(pRecv asymmetric.IPubKey, pPld payload.IPayload, pPadd uint64) (message.IMessage, error) {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.GetBytes(symmetric.CAESKeySize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	payloadBytes := pPld.ToBytes()
	payloadWrapper := payload.NewPayload(
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
			p.GetPubKey().GetHasher().ToBytes(),
			pRecv.GetHasher().ToBytes(),
			payloadWrapper.ToBytes(),
		},
		[]byte{},
	)).ToBytes()

	encKey := pRecv.EncryptBytes(session)
	if encKey == nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewAESCipher(session)
	return &message.SMessage{
		FPubKey:  encoding.HexEncode(cipher.EncryptBytes(p.GetPubKey().ToBytes())),
		FEncKey:  encoding.HexEncode(encKey),
		FSalt:    encoding.HexEncode(cipher.EncryptBytes(salt)),
		FHash:    encoding.HexEncode(cipher.EncryptBytes(hash)),
		FSign:    encoding.HexEncode(cipher.EncryptBytes(p.fPrivKey.SignBytes(hash))),
		FPayload: cipher.EncryptBytes(payloadWrapper.ToBytes()), // JSON field to raw Body (no need HEX encode)
	}, nil
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pMsg message.IMessage) (asymmetric.IPubKey, payload.IPayload, error) {
	// Initial check.
	if pMsg == nil || !pMsg.IsValid(p.fSettings) {
		return nil, nil, ErrInitCheckMessage
	}

	// Decrypt session key by private key of receiver.
	session := p.fPrivKey.DecryptBytes(pMsg.GetEncKey())
	if session == nil {
		return nil, nil, ErrDecryptCipherKey
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := symmetric.NewAESCipher(session)
	publicBytes := cipher.DecryptBytes(pMsg.GetPubKey())

	// Load public key and check standart size.
	pubKey := asymmetric.LoadRSAPubKey(publicBytes)
	if pubKey == nil {
		return nil, nil, ErrDecryptPublicKey
	}
	if pubKey.GetSize() != p.GetPubKey().GetSize() {
		return nil, nil, ErrInvalidPublicKeySize
	}

	// Decrypt main data of message by session key.
	payloadWrapperBytes := cipher.DecryptBytes(pMsg.GetPayload())
	payloadWrapper := payload.LoadPayload(payloadWrapperBytes)
	if payloadWrapper == nil {
		return nil, nil, ErrDecodePayloadWrapper
	}

	// Check size of payload.
	mustLen := payloadWrapper.GetHead()
	payloadBytes := payloadWrapper.GetBody()
	if mustLen > uint64(len(payloadBytes)) {
		return nil, nil, ErrInvalidPayloadSize
	}

	// Decrypt salt & hash.
	salt := cipher.DecryptBytes(pMsg.GetSalt())
	hash := cipher.DecryptBytes(pMsg.GetHash())

	// Validate received hash with generated hash.
	check := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			pubKey.GetHasher().ToBytes(),
			p.GetPubKey().GetHasher().ToBytes(),
			payloadWrapper.ToBytes(),
		},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, hash) {
		return nil, nil, ErrInvalidDataHash
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.DecryptBytes(pMsg.GetSign())
	if !pubKey.VerifyBytes(hash, sign) {
		return nil, nil, ErrInvalidHashSign
	}

	// Remove random bytes and get main data
	pld := payload.LoadPayload(payloadBytes[:mustLen])
	if pld == nil {
		return nil, nil, ErrDecodePayload
	}

	// Return decrypted message with title
	return pubKey, pld, nil
}
