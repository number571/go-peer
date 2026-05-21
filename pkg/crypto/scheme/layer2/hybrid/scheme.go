package hybrid

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

const (
	cSaltSize = 32 // bytes
)

var (
	_ layer2.IScheme = &sScheme{}
)

// Basic structure describing the user.
type sScheme struct {
	fPrivKey      asymmetric.IPrivKey
	fMessageSize  uint64
	fPayloadLimit uint64
}

// Create scheme by private key as identification.
// Handle function is used when the network exists.
func NewScheme(pPrivKey asymmetric.IPrivKey, pMessageSize uint64) layer2.IScheme {
	scheme := &sScheme{
		fMessageSize: pMessageSize,
		fPrivKey:     pPrivKey,
	}

	pubKey := pPrivKey.GetPubKey()
	encMsg, err := scheme.encryptWithPadding(pubKey, []byte{}, 0)
	if err != nil {
		panic(err)
	}

	structSize := uint64(len(encMsg))
	if pMessageSize <= structSize {
		panic("the payload size is lower than struct size")
	}

	scheme.fPayloadLimit = pMessageSize - structSize
	return scheme
}

func (p *sScheme) GetRandomKey() layer2.IParticipantKey {
	return asymmetric.NewPrivKey().GetPubKey()
}

// Message is raw bytes of full structure+payload.
func (p *sScheme) GetMessageSize() uint64 {
	return p.fMessageSize
}

// Payload is raw bytes of message without structure.
func (p *sScheme) GetPayloadLimit() uint64 {
	return p.fPayloadLimit
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (p *sScheme) EncryptMessage(pRecv layer2.IParticipantKey, pMsg []byte) ([]byte, error) {
	recv, ok := pRecv.(asymmetric.IPubKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}

	var (
		payloadLimit = p.fPayloadLimit
		resultSize   = uint64(len(pMsg))
	)

	if resultSize > payloadLimit {
		return nil, ErrLimitMessageSize
	}

	return p.encryptWithPadding(recv, pMsg, payloadLimit-resultSize)
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the layer2.
func (p *sScheme) DecryptMessage(pMapPubKeys layer2.IKeysContainer, pMsg []byte) (layer2.IParticipantKey, []byte, error) {
	mapPubKeys, ok := pMapPubKeys.(asymmetric.IMapPubKeys)
	if !ok {
		return nil, nil, ErrInvalidKeyType
	}

	// Load message's structure from encrypted bytes.
	msg, err := loadMessage(p.fMessageSize, pMsg)
	if err != nil {
		return nil, nil, ErrInitCheckMessage
	}

	// Decrypt session key by private key of receiver.
	skey, err := p.fPrivKey.GetKEMPrivKey().Decapsulate(msg.GetEnck())
	if err != nil {
		return nil, nil, ErrDecryptCipherKey
	}

	// Decrypt data block by decrypted session key. Decode data block.
	decJoiner := symmetric.NewCipherCFB(skey).DecryptBytes(msg.GetEncd())
	decSlice, err := joiner.LoadBytesJoiner32(decJoiner)
	if err != nil || len(decSlice) != 5 {
		return nil, nil, ErrDecodeBytesJoiner
	}

	// Decode wrapped data.
	var (
		pkid = decSlice[0]
		salt = decSlice[1]
		data = decSlice[2]
		hash = decSlice[3]
		sign = decSlice[4]
	)

	// Get public key from map by pkid (hash).
	sPubKey := mapPubKeys.GetPubKey(pkid)
	if sPubKey == nil {
		return nil, nil, ErrDecodePublicKey
	}

	// Validate received hash with generated hash.
	check := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{sPubKey.ToBytes(), p.fPrivKey.GetPubKey().ToBytes(), data},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, hash) {
		return nil, nil, ErrInvalidDataHash
	}

	// Verify sign by public key of sender and hash of message.
	if !sPubKey.GetDSAPubKey().VerifyBytes(hash, sign) {
		return nil, nil, ErrInvalidHashSign
	}

	// Decode main data of message by session key.
	payloadWrapper, err := joiner.LoadBytesJoiner32(data)
	if err != nil || len(payloadWrapper) != 2 {
		return nil, nil, ErrDecodePayloadWrapper
	}

	// Return public key of sender with payload.
	return sPubKey, payloadWrapper[0], nil
}

func (p *sScheme) encryptWithPadding(
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

	cipher := symmetric.NewCipherCFB(sk)
	return newMessage(
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
