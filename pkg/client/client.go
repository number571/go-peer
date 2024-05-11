package client

import (
	"bytes"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/payload/joiner"
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
		payload.NewPayload64(0, []byte{}),
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
func (p *sClient) EncryptPayload(pRecv asymmetric.IPubKey, pPld payload.IPayload64) (message.IMessage, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pPld.GetBody()))
	)

	if resultSize > msgLimitSize {
		return nil, ErrLimitMessageSize
	}

	return p.encryptWithParams(
		pRecv,
		pPld,
		msgLimitSize-resultSize,
	)
}

func (p *sClient) encryptWithParams(
	pRecv asymmetric.IPubKey,
	pPld payload.IPayload64,
	pPadd uint64,
) (message.IMessage, error) {
	var (
		rand    = random.NewCSPRNG()
		salt    = rand.GetBytes(symmetric.CAESKeySize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	data := joiner.NewBytesJoiner32([][]byte{
		pPld.ToBytes(),
		rand.GetBytes(pPadd),
	})

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			p.GetPubKey().GetHasher().ToBytes(),
			pRecv.GetHasher().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	encKey := pRecv.EncryptBytes(session)
	if encKey == nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewAESCipher(session)
	return message.NewMessage(
		encKey,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			p.GetPubKey().ToBytes(),
			salt,
			hash,
			p.fPrivKey.SignBytes(hash),
			data,
		})),
	), nil
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pMsg message.IMessage) (asymmetric.IPubKey, payload.IPayload64, error) {
	if pMsg == nil || !pMsg.IsValid(p.fSettings) {
		return nil, nil, ErrInitCheckMessage
	}

	// Decrypt session key by private key of receiver.
	session := p.fPrivKey.DecryptBytes(pMsg.GetEnck())
	if session == nil {
		return nil, nil, ErrDecryptCipherKey
	}

	// Decrypt data block by decrypted session key. Decode data block.
	decJoiner := symmetric.NewAESCipher(session).DecryptBytes(pMsg.GetEncd())
	decSlice, err := joiner.LoadBytesJoiner32(decJoiner)
	if err != nil || len(decSlice) != 5 {
		return nil, nil, ErrDecodeBytesJoiner
	}

	// Decode wrapped data.
	var (
		pkey = decSlice[0]
		salt = decSlice[1]
		hash = decSlice[2]
		sign = decSlice[3]
		data = decSlice[4]
	)

	// Load public key and check standart size.
	pubKey := asymmetric.LoadRSAPubKey(pkey)
	if pubKey == nil || pubKey.GetSize() != p.GetPubKey().GetSize() {
		return nil, nil, ErrDecodePublicKey
	}

	// Validate received hash with generated hash.
	check := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			pubKey.GetHasher().ToBytes(),
			p.GetPubKey().GetHasher().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, hash) {
		return nil, nil, ErrInvalidDataHash
	}

	// Verify sign by public key of sender and hash of message.
	if !pubKey.VerifyBytes(hash, sign) {
		return nil, nil, ErrInvalidHashSign
	}

	// Decode main data of message by session key.
	payloadWrapper, err := joiner.LoadBytesJoiner32(data)
	if err != nil || len(payloadWrapper) != 2 {
		return nil, nil, ErrDecodePayloadWrapper
	}

	// Remove random bytes and get main data.
	pld := payload.LoadPayload64(payloadWrapper[0])
	if pld == nil {
		return nil, nil, ErrDecodePayload
	}

	// Return public key of sender with payload.
	return pubKey, pld, nil
}
