package client

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/message/layer2"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	cSaltSize = 32 // bytes
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fPrivKey             asymmetric.IPrivKey
	fMessageSize         uint64
	fStructSize          uint64
	fStaticDecryptValues *sStaticDecryptValues
}

type sStaticDecryptValues struct {
	fEncMsg         layer2.IMessage
	fSKey           []byte
	fDecSlice       [][]byte
	fPubKey         asymmetric.IPubKey
	fPayloadWrapper [][]byte
}

// Create client by private key as identification.
// Handle function is used when the network exists.
func NewClient(pPrivKey asymmetric.IPrivKey, pMessageSize uint64) IClient {
	client := &sClient{
		fMessageSize: pMessageSize,
		fPrivKey:     pPrivKey,
	}

	pubKey := client.GetPrivKey().GetPubKey()
	encMsg, err := client.encryptWithPadding(pubKey, []byte{}, 0)
	if err != nil {
		panic(err)
	}

	client.fStructSize = uint64(len(encMsg))
	if limit := client.GetPayloadLimit(); limit == 0 {
		panic("the payload size is lower than struct size")
	}

	return client.withStaticDecryptValues()
}

func (p *sClient) GetMessageSize() uint64 {
	return p.fMessageSize
}

// Message is raw bytes of body payload without payload head.
func (p *sClient) GetPayloadLimit() uint64 {
	maxMsgSize := p.fMessageSize
	if maxMsgSize <= p.fStructSize {
		return 0
	}
	return maxMsgSize - p.fStructSize
}

// Get private key from client object.
func (p *sClient) GetPrivKey() asymmetric.IPrivKey {
	return p.fPrivKey
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (p *sClient) EncryptMessage(pRecv asymmetric.IPubKey, pMsg []byte) ([]byte, error) {
	var (
		payloadLimit = p.GetPayloadLimit()
		resultSize   = uint64(len(pMsg))
	)

	if resultSize > payloadLimit {
		return nil, ErrLimitMessageSize
	}

	return p.encryptWithPadding(pRecv, pMsg, payloadLimit-resultSize)
}

func (p *sClient) encryptWithPadding(
	pRecv asymmetric.IPubKey,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
		pkey = p.fPrivKey.GetPubKey()
	)

	data := joiner.NewBytesJoiner32([][]byte{pMsg, rand.GetBytes(pPadd)})
	hash := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{pkey.ToBytes(), pRecv.ToBytes(), data},
		[]byte{},
	)).ToBytes()

	ct, sk, err := pRecv.GetKEMPubKey().Encapsulate()
	if err != nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewCipher(sk)
	return layer2.NewMessage(
		ct,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			pkey.GetHasher().ToBytes(),
			salt,
			data,
			hash,
			p.fPrivKey.GetDSAPrivKey().SignBytes(hash),
		})),
	).ToBytes(), nil
}

func (p *sClient) withStaticDecryptValues() *sClient {
	var (
		encMsg, _         = p.EncryptMessage(p.fPrivKey.GetPubKey(), []byte{})
		msg, _            = layer2.LoadMessage(p.fMessageSize, encMsg)
		skey, _           = p.fPrivKey.GetKEMPrivKey().Decapsulate(msg.GetEnck())
		decJoiner         = symmetric.NewCipher(skey).DecryptBytes(msg.GetEncd())
		decSlice, _       = joiner.LoadBytesJoiner32(decJoiner)
		pkid              = decSlice[0]
		data              = decSlice[2]
		mapPubKeys        = asymmetric.NewMapPubKeys(p.fPrivKey.GetPubKey())
		sPubKey           = mapPubKeys.GetPubKey(pkid)
		payloadWrapper, _ = joiner.LoadBytesJoiner32(data)
	)
	p.fStaticDecryptValues = &sStaticDecryptValues{
		fEncMsg:         msg,
		fSKey:           skey,
		fDecSlice:       decSlice,
		fPubKey:         sPubKey,
		fPayloadWrapper: payloadWrapper,
	}
	return p
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pMapPubKeys asymmetric.IMapPubKeys, pMsg []byte) (asymmetric.IPubKey, []byte, error) {
	var resultError error

	msg, err := layer2.LoadMessage(p.fMessageSize, pMsg)
	if err != nil {
		msg = p.fStaticDecryptValues.fEncMsg
		resultError = ErrInitCheckMessage
	}

	// Decrypt session key by private key of receiver.
	skey, err := p.fPrivKey.GetKEMPrivKey().Decapsulate(msg.GetEnck())
	if err != nil {
		skey = p.fStaticDecryptValues.fSKey
		if resultError == nil {
			resultError = ErrDecryptCipherKey
		}
	}

	// Decrypt data block by decrypted session key. Decode data block.
	decJoiner := symmetric.NewCipher(skey).DecryptBytes(msg.GetEncd())
	decSlice, err := joiner.LoadBytesJoiner32(decJoiner)
	if err != nil || len(decSlice) != 5 {
		decSlice = p.fStaticDecryptValues.fDecSlice
		if resultError == nil {
			resultError = ErrDecodeBytesJoiner
		}
	}

	// Decode wrapped data.
	var (
		pkid = decSlice[0]
		salt = decSlice[1]
		data = decSlice[2]
		hash = decSlice[3]
		sign = decSlice[4]
	)

	// Get public key from map by pkid (hash)
	sPubKey := pMapPubKeys.GetPubKey(pkid)
	if sPubKey == nil {
		sPubKey = p.fStaticDecryptValues.fPubKey
		if resultError == nil {
			resultError = ErrDecodePublicKey
		}
	}

	// Validate received hash with generated hash.
	check := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{sPubKey.ToBytes(), p.fPrivKey.GetPubKey().ToBytes(), data},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, hash) {
		if resultError == nil {
			resultError = ErrInvalidDataHash
		}
	}

	// Verify sign by public key of sender and hash of message.
	if !sPubKey.GetDSAPubKey().VerifyBytes(hash, sign) {
		if resultError == nil {
			resultError = ErrInvalidHashSign
		}
	}

	// Decode main data of message by session key.
	payloadWrapper, err := joiner.LoadBytesJoiner32(data)
	if err != nil || len(payloadWrapper) != 2 {
		payloadWrapper = p.fStaticDecryptValues.fPayloadWrapper
		if resultError == nil {
			resultError = ErrDecodePayloadWrapper
		}
	}

	// Return error if exists
	if err := resultError; err != nil {
		return nil, nil, err
	}

	// Return public key of sender with payload.
	return sPubKey, payloadWrapper[0], nil
}
