package local

import (
	"bytes"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/settings"
)

var (
	_ Client = &clientT{}
)

// Basic structure describing the user.
type clientT struct {
	gs settings.Settings
	pk crypto.PrivKey
}

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv crypto.PrivKey, s settings.Settings) Client {
	if priv == nil {
		return nil
	}
	return &clientT{
		gs: s,
		pk: priv,
	}
}

// Get public key from client object.
func (client *clientT) PubKey() crypto.PubKey {
	return client.pk.PubKey()
}

// Get private key from client object.
func (client *clientT) PrivKey() crypto.PrivKey {
	return client.pk
}

// Get settings from client object.
func (client *clientT) Settings() settings.Settings {
	return client.gs
}

// Function wrap message in multiple route encryption.
// Need use pseudo sender if route not null.
func (client *clientT) Encrypt(route Route, msg Message) (Message, Session) {
	var (
		psender       = NewClient(route.psender, client.gs)
		rmsg, session = client.onceEncrypt(route.receiver, msg)
	)
	if psender == nil && len(route.routes) != 0 {
		return nil, nil
	}
	for _, pub := range route.routes {
		rmsg, _ = psender.(*clientT).onceEncrypt(
			pub,
			NewMessage(
				encoding.Uint64ToBytes(client.gs.Get(settings.MaskRout)),
				rmsg.Serialize().Bytes(),
			),
		)
	}
	return rmsg, session
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *clientT) onceEncrypt(receiver crypto.PubKey, msg Message) (Message, Session) {
	var (
		session   = crypto.RandBytes(client.gs.Get(settings.SizeSkey))
		randBytes = crypto.RandBytes(client.gs.Get(settings.SizeSkey))
		cipher    = crypto.NewCipher(session)
	)

	data := bytes.Join(
		[][]byte{
			encoding.Uint64ToBytes(uint64(len(msg.Body.Data))),
			msg.Body.Data,
			encoding.Uint64ToBytes(crypto.RandUint64() % (settings.SizePack / 4)),
		},
		[]byte{},
	)

	hash := crypto.NewHasher(bytes.Join(
		[][]byte{
			randBytes,
			client.PubKey().Bytes(),
			receiver.Bytes(),
			data,
		},
		[]byte{},
	)).Bytes()

	return &messageT{
		Head: headMessage{
			Sender:    cipher.Encrypt(client.PubKey().Bytes()),
			Session:   receiver.Encrypt(session),
			RandBytes: cipher.Encrypt(randBytes),
		},
		Body: bodyMessage{
			Data:  cipher.Encrypt(data),
			Hash:  hash,
			Sign:  cipher.Encrypt(client.PrivKey().Sign(hash)),
			Proof: crypto.NewPuzzle(client.Settings().Get(settings.SizeWork)).Proof(hash),
		},
	}, session
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *clientT) Decrypt(msg Message) Message {
	const (
		SizeUint64 = 8 // bytes
	)

	// Initial check.
	if len(msg.Body.Hash) != crypto.HashSize {
		return nil
	}

	// Proof of work. Prevent spam.
	diff := client.Settings().Get(settings.SizeWork)
	puzzle := crypto.NewPuzzle(diff)
	if !puzzle.Verify(msg.Body.Hash, msg.Body.Proof) {
		return nil
	}

	// Decrypt session key by private key of receiver.
	session := client.PrivKey().Decrypt(msg.Head.Session)
	if session == nil {
		return nil
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := crypto.NewCipher(session)
	publicBytes := cipher.Decrypt(msg.Head.Sender)
	if publicBytes == nil {
		return nil
	}

	// Load public key and check standart size.
	public := crypto.LoadPubKey(publicBytes)
	if public == nil {
		return nil
	}
	if public.Size() != client.PubKey().Size() {
		return nil
	}

	// Decrypt random bytes.
	randBytes := cipher.Decrypt(msg.Head.RandBytes)
	if randBytes == nil {
		return nil
	}

	// Decrypt main data of message by session key.
	dataBytes := cipher.Decrypt(msg.Body.Data)
	if dataBytes == nil {
		return nil
	}

	// Check received hash and generated hash.
	check := crypto.NewHasher(bytes.Join(
		[][]byte{
			randBytes,
			publicBytes,
			client.PubKey().Bytes(),
			dataBytes,
		},
		[]byte{},
	)).Bytes()
	if !bytes.Equal(check, msg.Body.Hash) {
		return nil
	}

	// check size of data
	if len(dataBytes) < SizeUint64 {
		return nil
	}

	// pass random bytes and get main data
	mustLen := encoding.BytesToUint64(dataBytes[:SizeUint64])
	allData := dataBytes[SizeUint64:]
	if mustLen > uint64(len(allData)) {
		return nil
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.Decrypt(msg.Body.Sign)
	if sign == nil {
		return nil
	}
	if !public.Verify(msg.Body.Hash, sign) {
		return nil
	}

	// Return decrypted message.
	return &messageT{
		Head: headMessage{
			Sender:    publicBytes,
			Session:   session,
			RandBytes: randBytes,
		},
		Body: bodyMessage{
			Data:  allData[:mustLen],
			Hash:  msg.Body.Hash,
			Sign:  sign,
			Proof: msg.Body.Proof,
		},
	}
}
