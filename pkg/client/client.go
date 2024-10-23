package client

import (
	"bytes"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
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
	fPrivKey     asymmetric.IPrivKey
	fMessageSize uint64
	fStructSize  uint64
}

// Create client by private key as identification.
// Handle function is used when the network exists.
func NewClient(pPrivKey asymmetric.IPrivKey, pMessageSize uint64) IClient {
	client := &sClient{
		fMessageSize: pMessageSize,
		fPrivKey:     pPrivKey,
	}

	kemPubKey := client.GetPrivKey().GetKEMPrivKey().GetPubKey()
	encMsg, err := client.encryptWithParams(kemPubKey, []byte{}, 0)
	if err != nil {
		panic(err)
	}

	client.fStructSize = uint64(len(encMsg))
	if limit := client.GetPayloadLimit(); limit == 0 {
		panic("the payload size is lower than struct size")
	}

	return client
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
func (p *sClient) EncryptMessage(pRecv asymmetric.IKEMPubKey, pMsg []byte) ([]byte, error) {
	var (
		payloadLimit = p.GetPayloadLimit()
		resultSize   = uint64(len(pMsg))
	)

	if resultSize > payloadLimit {
		return nil, ErrLimitMessageSize
	}

	return p.encryptWithParams(pRecv, pMsg, payloadLimit-resultSize)
}

func (p *sClient) encryptWithParams(
	pRecv asymmetric.IKEMPubKey,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		rand = random.NewRandom()
		salt = rand.GetBytes(cSaltSize)
		pkey = p.fPrivKey.GetPubKey().ToBytes()
	)

	data := joiner.NewBytesJoiner32([][]byte{pMsg, rand.GetBytes(pPadd)})
	hash := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			pkey,
			pRecv.ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()

	ct, sk, err := pRecv.Encapsulate()
	if err != nil {
		return nil, ErrEncryptSymmetricKey
	}

	cipher := symmetric.NewCipher(sk)
	return message.NewMessage(
		ct,
		cipher.EncryptBytes(joiner.NewBytesJoiner32([][]byte{
			pkey,
			salt,
			hash,
			p.fPrivKey.GetDSAPrivKey().SignBytes(hash),
			data,
		})),
	).ToBytes(), nil
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pMsg []byte) (asymmetric.IPubKey, []byte, error) {
	msg, err := message.LoadMessage(p.fMessageSize, pMsg)
	if err != nil {
		return nil, nil, ErrInitCheckMessage
	}

	// Decrypt session key by private key of receiver.
	kemPrivKey := p.fPrivKey.GetKEMPrivKey()
	skey, err := kemPrivKey.Decapsulate(msg.GetEnck())
	if err != nil {
		return nil, nil, ErrDecryptCipherKey
	}

	// Decrypt data block by decrypted session key. Decode data block.
	decJoiner := symmetric.NewCipher(skey).DecryptBytes(msg.GetEncd())
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
	pubKey := asymmetric.LoadPubKey(pkey)
	if pubKey == nil {
		return nil, nil, ErrDecodePublicKey
	}

	// Validate received hash with generated hash.
	check := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			pubKey.ToBytes(),
			kemPrivKey.GetPubKey().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, hash) {
		return nil, nil, ErrInvalidDataHash
	}

	// Verify sign by public key of sender and hash of message.
	if !pubKey.GetDSAPubKey().VerifyBytes(hash, sign) {
		return nil, nil, ErrInvalidHashSign
	}

	// Decode main data of message by session key.
	payloadWrapper, err := joiner.LoadBytesJoiner32(data)
	if err != nil || len(payloadWrapper) != 2 {
		return nil, nil, ErrDecodePayloadWrapper
	}

	// Return public key of sender with payload.
	return pubKey, payloadWrapper[0], nil
}
