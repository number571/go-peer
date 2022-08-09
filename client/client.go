package client

import (
	"bytes"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/puzzle"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/crypto/symmetric"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/routing"
	"github.com/number571/go-peer/settings"
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fSettings ISettings
	fPrivKey  asymmetric.IPrivKey
}

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(sett ISettings, priv asymmetric.IPrivKey) IClient {
	if priv == nil {
		return nil
	}
	return &sClient{
		fSettings: sett,
		fPrivKey:  priv,
	}
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

// Function wrap message in multiple route encryption.
// Need use pseudo sender if route not null.
func (client *sClient) Encrypt(route routing.IRoute, pl payload.IPayload) message.IMessage {
	var (
		psender = NewClient(client.Settings(), route.PSender())
		rmsg    = client.onceEncrypt(route.Receiver(), pl)
	)
	if psender == nil && len(route.List()) != 0 {
		return nil
	}
	for _, pub := range route.List() {
		rmsg = psender.(*sClient).onceEncrypt(
			pub,
			payload.NewPayload(
				settings.CMaskRoute,
				rmsg.Bytes(),
			),
		)
	}
	return rmsg
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *sClient) onceEncrypt(receiver asymmetric.IPubKey, pl payload.IPayload) message.IMessage {
	var (
		rand    = random.NewStdPRNG()
		salt    = rand.Bytes(settings.CSizeSymmKey)
		session = rand.Bytes(settings.CSizeSymmKey)
	)

	maxRandSize := client.Settings().GetRandomSize()
	if maxRandSize == 0 {
		maxRandSize = 1
	}

	randBytes := rand.Bytes(rand.Uint64() % maxRandSize)
	doublePayload := payload.NewPayload(
		uint64(len(pl.Bytes())), // head as size of (payload||random)
		bytes.Join(
			[][]byte{
				pl.Bytes(),
				randBytes,
			},
			[]byte{},
		),
	)

	hash := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			salt,
			client.PubKey().Bytes(),
			receiver.Bytes(),
			doublePayload.Bytes(),
		},
		[]byte{},
	)).Bytes()

	cipher := symmetric.NewAESCipher(session)
	return &message.SMessage{
		FHead: message.SHeadMessage{
			FSender:  cipher.Encrypt(client.PubKey().Bytes()),
			FSession: receiver.Encrypt(session),
			FSalt:    cipher.Encrypt(salt),
		},
		FBody: message.SBodyMessage{
			FPayload: cipher.Encrypt(doublePayload.Bytes()),
			FHash:    hash,
			FSign:    cipher.Encrypt(client.PrivKey().Sign(hash)),
			FProof:   puzzle.NewPoWPuzzle(client.Settings().GetWorkSize()).Proof(hash),
		},
	}
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *sClient) Decrypt(msg message.IMessage) (asymmetric.IPubKey, payload.IPayload) {
	if msg == nil {
		return nil, nil
	}

	// Initial check.
	if len(msg.Body().Hash()) != hashing.GSHA256Size {
		return nil, nil
	}

	// Proof of work. Prevent spam.
	diff := client.Settings().GetWorkSize()
	puzzle := puzzle.NewPoWPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil, nil
	}

	// Decrypt session key by private key of receiver.
	session := client.PrivKey().Decrypt(msg.Head().Session())
	if session == nil {
		return nil, nil
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := symmetric.NewAESCipher(session)
	publicBytes := cipher.Decrypt(msg.Head().Sender())
	if publicBytes == nil {
		return nil, nil
	}

	// Load public key and check standart size.
	pubKey := asymmetric.LoadRSAPubKey(publicBytes)
	if pubKey == nil {
		return nil, nil
	}
	if pubKey.Size() != client.PubKey().Size() {
		return nil, nil
	}

	// Decrypt salt.
	salt := cipher.Decrypt(msg.Head().Salt())
	if salt == nil {
		return nil, nil
	}

	// Decrypt main data of message by session key.
	doublePayloadBytes := cipher.Decrypt(msg.Body().Payload().Bytes())
	if doublePayloadBytes == nil {
		return nil, nil
	}
	doublePayload := payload.LoadPayload(doublePayloadBytes)
	if doublePayload == nil {
		return nil, nil
	}

	// Check received hash and generated hash.
	check := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			salt,
			publicBytes,
			client.PubKey().Bytes(),
			doublePayload.Bytes(),
		},
		[]byte{},
	)).Bytes()
	if !bytes.Equal(check, msg.Body().Hash()) {
		return nil, nil
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.Decrypt(msg.Body().Sign())
	if sign == nil {
		return nil, nil
	}
	if !pubKey.Verify(msg.Body().Hash(), sign) {
		return nil, nil
	}

	// remove random bytes and get main data
	mustLen := doublePayload.Head()
	if mustLen > uint64(len(doublePayload.Body())) {
		return nil, nil
	}
	payload := payload.LoadPayload(doublePayload.Body()[:mustLen])
	if payload == nil {
		return nil, nil
	}

	// Return decrypted message with title
	return pubKey, payload
}
