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
	fSettings     message.ISettings
	fPrivKeyChain asymmetric.IPrivKeyChain
	fStructSize   uint64
}

// Create client by private key as identification.
// Handle function is used when the network exists.
func NewClient(pSett message.ISettings, pPrivKeyChain asymmetric.IPrivKeyChain) IClient {
	client := &sClient{
		fSettings:     pSett,
		fPrivKeyChain: pPrivKeyChain,
	}

	kemPubKey := client.GetPrivKeyChain().GetKEncPrivKey().GetPubKey()
	encMsg, err := client.encryptWithParams(kemPubKey, []byte{}, 0)
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

// Get private key from client object.
func (p *sClient) GetPrivKeyChain() asymmetric.IPrivKeyChain {
	return p.fPrivKeyChain
}

// Get settings from client object.
func (p *sClient) GetSettings() message.ISettings {
	return p.fSettings
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (p *sClient) EncryptMessage(pRecv asymmetric.IKEncPubKey, pMsg []byte) ([]byte, error) {
	var (
		msgLimitSize = p.GetMessageLimit()
		resultSize   = uint64(len(pMsg))
	)

	if resultSize > msgLimitSize {
		return nil, ErrLimitMessageSize
	}

	return p.encryptWithParams(pRecv, pMsg, msgLimitSize-resultSize)
}

func (p *sClient) encryptWithParams(
	pRecv asymmetric.IKEncPubKey,
	pMsg []byte,
	pPadd uint64,
) ([]byte, error) {
	var (
		signer = p.fPrivKeyChain.GetSignPrivKey()
		rand   = random.NewRandom()
		salt   = rand.GetBytes(cSaltSize)
	)

	data := joiner.NewBytesJoiner32([][]byte{pMsg, rand.GetBytes(pPadd)})
	hash := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			signer.GetPubKey().ToBytes(),
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
			signer.GetPubKey().ToBytes(),
			salt,
			hash,
			signer.SignBytes(hash),
			data,
		})),
	).ToBytes(), nil
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (p *sClient) DecryptMessage(pMsg []byte) (asymmetric.ISignPubKey, []byte, error) {
	msg, err := message.LoadMessage(p.fSettings, pMsg)
	if err != nil {
		return nil, nil, ErrInitCheckMessage
	}

	// Decrypt session key by private key of receiver.
	kemPrivKey := p.fPrivKeyChain.GetKEncPrivKey()
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
	signerPubKey := asymmetric.LoadSignPubKey(pkey)
	if signerPubKey == nil { // || pubKey.GetSize() != p.GetPubKey().GetSize()
		return nil, nil, ErrDecodePublicKey
	}

	// Validate received hash with generated hash.
	check := hashing.NewHMACHasher(salt, bytes.Join(
		[][]byte{
			signerPubKey.ToBytes(),
			kemPrivKey.GetPubKey().ToBytes(),
			data,
		},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, hash) {
		return nil, nil, ErrInvalidDataHash
	}

	// Verify sign by public key of sender and hash of message.
	if !signerPubKey.VerifyBytes(hash, sign) {
		return nil, nil, ErrInvalidHashSign
	}

	// Decode main data of message by session key.
	payloadWrapper, err := joiner.LoadBytesJoiner32(data)
	if err != nil || len(payloadWrapper) != 2 {
		return nil, nil, ErrDecodePayloadWrapper
	}

	// Return public key of sender with payload.
	return signerPubKey, payloadWrapper[0], nil
}
