package client

import (
	"bytes"
	"fmt"

	"github.com/number571/go-peer/modules/client/message"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/crypto/puzzle"
	"github.com/number571/go-peer/modules/crypto/random"
	"github.com/number571/go-peer/modules/crypto/symmetric"
	"github.com/number571/go-peer/modules/encoding"
	"github.com/number571/go-peer/modules/payload"
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
// Handle function is used when the network exists. Can be null.
func NewClient(sett ISettings, priv asymmetric.IPrivKey) IClient {
	if priv == nil {
		return nil
	}
	client := &sClient{
		fSettings: sett,
		fPrivKey:  priv,
	}
	msg := client.encryptWithParams(
		client.PubKey(),
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
func (client *sClient) PubKey() asymmetric.IPubKey {
	return client.PrivKey().PubKey()
}

// Get private key from client object.
func (client *sClient) PrivKey() asymmetric.IPrivKey {
	return client.fPrivKey
}

// Get settings from client object.
func (client *sClient) Settings() ISettings {
	return client.fSettings
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *sClient) Encrypt(receiver asymmetric.IPubKey, pl payload.IPayload) (message.IMessage, error) {
	if receiver.Size() != client.PubKey().Size() {
		return nil, fmt.Errorf("size of public keys sender and receiver not equal")
	}

	var (
		maxMsgSize = client.Settings().GetMessageSize() >> 1 // limit of bytes without hex
		resultSize = uint64(client.fVoidMsgSize) + uint64(len(pl.ToBytes()))
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
		pl,
		client.Settings().GetWorkSize(),
		maxMsgSize-resultSize,
	), nil
}

func (client *sClient) encryptWithParams(receiver asymmetric.IPubKey, pl payload.IPayload, workSize, addPadd uint64) message.IMessage {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.Bytes(symmetric.CAESKeySize)
		session = rand.Bytes(symmetric.CAESKeySize)
	)

	payloadBytes := pl.ToBytes()
	doublePayload := payload.NewPayload(
		uint64(len(payloadBytes)),
		bytes.Join(
			[][]byte{
				payloadBytes,
				rand.Bytes(addPadd),
			},
			[]byte{},
		),
	)

	hash := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			salt,
			client.PubKey().Bytes(),
			receiver.Bytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).Bytes()

	cipher := symmetric.NewAESCipher(session)
	bProof := encoding.Uint64ToBytes(puzzle.NewPoWPuzzle(workSize).Proof(hash))
	return &message.SMessage{
		FHead: message.SHeadMessage{
			FSender:  encoding.HexEncode(cipher.Encrypt(client.PubKey().Bytes())),
			FSession: encoding.HexEncode(receiver.Encrypt(session)),
			FSalt:    encoding.HexEncode(cipher.Encrypt(salt)),
		},
		FBody: message.SBodyMessage{
			FPayload: encoding.HexEncode(cipher.Encrypt(doublePayload.ToBytes())),
			FHash:    encoding.HexEncode(hash),
			FSign:    encoding.HexEncode(cipher.Encrypt(client.PrivKey().Sign(hash))),
			FProof:   encoding.HexEncode(bProof[:]),
		},
	}
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *sClient) Decrypt(msg message.IMessage) (asymmetric.IPubKey, payload.IPayload, error) {
	if msg == nil {
		return nil, nil, fmt.Errorf("msg is nil")
	}

	// Initial check.
	if len(msg.Body().Hash()) != hashing.CSHA256Size {
		return nil, nil, fmt.Errorf("msg hash != sha256 size")
	}

	// Proof of work. Prevent spam.
	diff := client.Settings().GetWorkSize()
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil, nil, fmt.Errorf("invalid proof of msg")
	}

	// Decrypt session key by private key of receiver.
	session := client.PrivKey().Decrypt(msg.Head().Session())
	if session == nil {
		return nil, nil, fmt.Errorf("failed decrypt session key")
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := symmetric.NewAESCipher(session)
	publicBytes := cipher.Decrypt(msg.Head().Sender())
	if publicBytes == nil {
		return nil, nil, fmt.Errorf("failed decrypt public key")
	}

	// Load public key and check standart size.
	pubKey := asymmetric.LoadRSAPubKey(publicBytes)
	if pubKey == nil {
		return nil, nil, fmt.Errorf("failed load public key")
	}
	if pubKey.Size() != client.PubKey().Size() {
		return nil, nil, fmt.Errorf("invalid public key size")
	}

	// Decrypt salt.
	salt := cipher.Decrypt(msg.Head().Salt())
	if salt == nil {
		return nil, nil, fmt.Errorf("failed decrypt salt")
	}

	// Decrypt main data of message by session key.
	doublePayloadBytes := cipher.Decrypt(msg.Body().Payload().ToBytes())
	if doublePayloadBytes == nil {
		return nil, nil, fmt.Errorf("failed decrypt double payload")
	}
	doublePayload := payload.LoadPayload(doublePayloadBytes)
	if doublePayload == nil {
		return nil, nil, fmt.Errorf("failed load double payload")
	}

	// Check received hash and generated hash.
	check := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			salt,
			publicBytes,
			client.PubKey().Bytes(),
			doublePayload.ToBytes(),
		},
		[]byte{},
	)).Bytes()
	if !bytes.Equal(check, msg.Body().Hash()) {
		return nil, nil, fmt.Errorf("invalid msg hash")
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.Decrypt(msg.Body().Sign())
	if sign == nil {
		return nil, nil, fmt.Errorf("failed decrypt sign")
	}
	if !pubKey.Verify(msg.Body().Hash(), sign) {
		return nil, nil, fmt.Errorf("invalid msg sign")
	}

	// remove random bytes and get main data
	mustLen := doublePayload.Head()
	if mustLen > uint64(len(doublePayload.Body())) {
		return nil, nil, fmt.Errorf("invalid size of payload")
	}
	payload := payload.LoadPayload(doublePayload.Body()[:mustLen])
	if payload == nil {
		return nil, nil, fmt.Errorf("invalid load payload")
	}

	// Return decrypted message with title
	return pubKey, payload, nil
}
