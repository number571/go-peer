package client

import (
	"bytes"
	"fmt"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fSettings    ISettings
	fPrivKey     asymmetric.IPrivKey
	fVoidMsgSize int
}

// Create client by private key as identification.
// Handle function is used when the network exists.
func NewClient(sett ISettings, priv asymmetric.IPrivKey) IClient {
	client := &sClient{
		fSettings: sett,
		fPrivKey:  priv,
	}
	msg := client.encryptWithParams(
		client.GetPubKey(),
		payload.NewPayload(0, []byte{}),
		0,
		0,
	)
	// saved message size with hex encoding
	// because exists not encoded chars <{}",>
	// of JSON format
	client.fVoidMsgSize = len(msg.ToBytes())
	return client
}

// Get public key from client object.
func (client *sClient) GetPubKey() asymmetric.IPubKey {
	return client.GetPrivKey().PubKey()
}

// Get private key from client object.
func (client *sClient) GetPrivKey() asymmetric.IPrivKey {
	return client.fPrivKey
}

// Get settings from client object.
func (client *sClient) GetSettings() ISettings {
	return client.fSettings
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *sClient) EncryptPayload(receiver asymmetric.IPubKey, pld payload.IPayload) (message.IMessage, error) {
	if receiver.GetSize() != client.GetPubKey().GetSize() {
		return nil, fmt.Errorf("size of public keys sender and receiver not equal")
	}

	var (
		maxMsgSize = client.GetSettings().GetMessageSize() >> 1 // limit of bytes without hex
		resultSize = uint64(client.fVoidMsgSize) + uint64(len(pld.ToBytes()))
	)

	if resultSize > maxMsgSize {
		return nil, fmt.Errorf(
			"limit of message size without hex encoding = %d bytes < current payload size with additional padding = %d bytes",
			maxMsgSize,
			resultSize,
		)
	}

	return client.encryptWithParams(
		receiver,
		pld,
		client.GetSettings().GetWorkSize(),
		maxMsgSize-resultSize,
	), nil
}

func (client *sClient) encryptWithParams(receiver asymmetric.IPubKey, pld payload.IPayload, workSize, addPadd uint64) message.IMessage {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.GetBytes(symmetric.CAESKeySize)
		session = rand.GetBytes(symmetric.CAESKeySize)
	)

	payloadBytes := pld.ToBytes()
	doublePayload := payload.NewPayload(
		uint64(len(payloadBytes)),
		bytes.Join(
			[][]byte{
				payloadBytes,
				rand.GetBytes(addPadd),
			},
			[]byte{},
		),
	)

	hash := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			client.GetPubKey().Address().ToBytes(),
			receiver.Address().ToBytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).ToBytes()

	cipher := symmetric.NewAESCipher(session)
	bProof := encoding.Uint64ToBytes(puzzle.NewPoWPuzzle(workSize).Proof(hash))
	return &message.SMessage{
		FHead: message.SHeadMessage{
			FSender:  encoding.HexEncode(cipher.EncryptBytes(client.GetPubKey().ToBytes())),
			FSession: encoding.HexEncode(receiver.EncryptBytes(session)),
			FSalt:    encoding.HexEncode(cipher.EncryptBytes(salt)),
		},
		FBody: message.SBodyMessage{
			FPayload: encoding.HexEncode(cipher.EncryptBytes(doublePayload.ToBytes())),
			FHash:    encoding.HexEncode(hash),
			FSign:    encoding.HexEncode(cipher.EncryptBytes(client.GetPrivKey().Sign(hash))),
			FProof:   encoding.HexEncode(bProof[:]),
		},
	}
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *sClient) DecryptMessage(msg message.IMessage) (asymmetric.IPubKey, payload.IPayload, error) {
	if msg == nil {
		return nil, nil, fmt.Errorf("msg is nil")
	}

	// Initial check.
	if len(msg.GetBody().GetHash()) != hashing.CSHA256Size {
		return nil, nil, fmt.Errorf("msg hash != sha256 size")
	}

	// Proof of work. Prevent spam.
	diff := client.GetSettings().GetWorkSize()
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.GetBody().GetHash(), msg.GetBody().GetProof()) {
		return nil, nil, fmt.Errorf("invalid proof of msg")
	}

	// Decrypt session key by private key of receiver.
	session := client.GetPrivKey().DecryptBytes(msg.GetHead().GetSession())
	if session == nil {
		return nil, nil, fmt.Errorf("failed decrypt session key")
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := symmetric.NewAESCipher(session)
	publicBytes := cipher.DecryptBytes(msg.GetHead().GetSender())
	if publicBytes == nil {
		return nil, nil, fmt.Errorf("failed decrypt public key")
	}

	// Load public key and check standart size.
	pubKey := asymmetric.LoadRSAPubKey(publicBytes)
	if pubKey == nil {
		return nil, nil, fmt.Errorf("failed load public key")
	}
	if pubKey.GetSize() != client.GetPubKey().GetSize() {
		return nil, nil, fmt.Errorf("invalid public key size")
	}

	// Decrypt main data of message by session key.
	doublePayloadBytes := cipher.DecryptBytes(msg.GetBody().GetPayload().ToBytes())
	if doublePayloadBytes == nil {
		return nil, nil, fmt.Errorf("failed decrypt double payload")
	}
	doublePayload := payload.LoadPayload(doublePayloadBytes)
	if doublePayload == nil {
		return nil, nil, fmt.Errorf("failed load double payload")
	}

	// Decrypt salt.
	salt := cipher.DecryptBytes(msg.GetHead().GetSalt())
	if salt == nil {
		return nil, nil, fmt.Errorf("failed decrypt salt")
	}

	// Check received hash and generated hash.
	check := hashing.NewHMACSHA256Hasher(salt, bytes.Join(
		[][]byte{
			pubKey.Address().ToBytes(),
			client.GetPubKey().Address().ToBytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).ToBytes()
	if !bytes.Equal(check, msg.GetBody().GetHash()) {
		return nil, nil, fmt.Errorf("invalid msg hash")
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.DecryptBytes(msg.GetBody().GetSign())
	if sign == nil {
		return nil, nil, fmt.Errorf("failed decrypt sign")
	}
	if !pubKey.Verify(msg.GetBody().GetHash(), sign) {
		return nil, nil, fmt.Errorf("invalid msg sign")
	}

	// Remove random bytes and get main data
	mustLen := doublePayload.GetHead()
	if mustLen > uint64(len(doublePayload.GetBody())) {
		return nil, nil, fmt.Errorf("invalid size of payload")
	}
	pld := payload.LoadPayload(doublePayload.GetBody()[:mustLen])
	if pld == nil {
		return nil, nil, fmt.Errorf("invalid load payload")
	}

	// Return decrypted message with title
	return pubKey, pld, nil
}
