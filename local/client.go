package local

import (
	"bytes"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/settings"
)

var (
	_ IClient = &sClient{}
)

// Basic structure describing the user.
type sClient struct {
	fSettings settings.ISettings
	fPrivKey  crypto.IPrivKey
}

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv crypto.IPrivKey, sett settings.ISettings) IClient {
	if priv == nil {
		return nil
	}
	return &sClient{
		fSettings: sett,
		fPrivKey:  priv,
	}
}

// Get public key from client object.
func (client *sClient) PubKey() crypto.IPubKey {
	return client.PrivKey().PubKey()
}

// Get private key from client object.
func (client *sClient) PrivKey() crypto.IPrivKey {
	return client.fPrivKey
}

// Get settings from client object.
func (client *sClient) Settings() settings.ISettings {
	return client.fSettings
}

// Function wrap message in multiple route encryption.
// Need use pseudo sender if route not null.
func (client *sClient) Encrypt(route IRoute, msg IMessage) (IMessage, Session) {
	var (
		psender       = NewClient(route.PSender(), client.Settings())
		rmsg, session = client.onceEncrypt(route.Receiver(), msg)
	)
	if psender == nil && len(route.List()) != 0 {
		return nil, nil
	}
	for _, pub := range route.List() {
		rmsg, _ = psender.(*sClient).onceEncrypt(
			pub,
			NewMessage(
				encoding.Uint64ToBytes(client.Settings().Get(settings.MaskRout)),
				rmsg.ToPackage().Bytes(),
			),
		)
	}
	return rmsg, session
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *sClient) onceEncrypt(receiver crypto.IPubKey, msg IMessage) (IMessage, []byte) {
	var (
		rand    = crypto.NewPRNG()
		salt    = rand.Bytes(client.Settings().Get(settings.SizeSkey))
		session = rand.Bytes(client.Settings().Get(settings.SizeSkey))
	)

	data := bytes.Join(
		[][]byte{
			encoding.Uint64ToBytes(uint64(len(msg.Body().Data()))),
			msg.Body().Data(),
			encoding.Uint64ToBytes(rand.Uint64() % (settings.SizePack / 4)),
		},
		[]byte{},
	)

	hash := crypto.NewHasher(bytes.Join(
		[][]byte{
			salt,
			client.PubKey().Bytes(),
			receiver.Bytes(),
			data,
		},
		[]byte{},
	)).Bytes()

	cipher := crypto.NewCipher(session)
	return &sMessage{
		FHead: sHeadMessage{
			FSender:  cipher.Encrypt(client.PubKey().Bytes()),
			FSession: receiver.Encrypt(session),
			FSalt:    cipher.Encrypt(salt),
		},
		FBody: sBodyMessage{
			FData:  cipher.Encrypt(data),
			FHash:  hash,
			FSign:  cipher.Encrypt(client.PrivKey().Sign(hash)),
			FProof: crypto.NewPuzzle(client.Settings().Get(settings.SizeWork)).Proof(hash),
		},
	}, session
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *sClient) Decrypt(msg IMessage) (IMessage, Title) {
	const (
		SizeUint64 = 8 // bytes
	)

	if msg == nil {
		return nil, nil
	}

	// Initial check.
	if len(msg.Body().Hash()) != crypto.HashSize {
		return nil, nil
	}

	// Proof of work. Prevent spam.
	diff := client.Settings().Get(settings.SizeWork)
	puzzle := crypto.NewPuzzle(diff)
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		return nil, nil
	}

	// Decrypt session key by private key of receiver.
	session := client.PrivKey().Decrypt(msg.Head().Session())
	if session == nil {
		return nil, nil
	}

	// Decrypt public key of sender by decrypted session key.
	cipher := crypto.NewCipher(session)
	publicBytes := cipher.Decrypt(msg.Head().Sender())
	if publicBytes == nil {
		return nil, nil
	}

	// Load public key and check standart size.
	public := crypto.LoadPubKey(publicBytes)
	if public == nil {
		return nil, nil
	}
	if public.Size() != client.PubKey().Size() {
		return nil, nil
	}

	// Decrypt salt.
	salt := cipher.Decrypt(msg.Head().Salt())
	if salt == nil {
		return nil, nil
	}

	// Decrypt main data of message by session key.
	dataBytes := cipher.Decrypt(msg.Body().Data())
	if dataBytes == nil {
		return nil, nil
	}

	// Check received hash and generated hash.
	check := crypto.NewHasher(bytes.Join(
		[][]byte{
			salt,
			publicBytes,
			client.PubKey().Bytes(),
			dataBytes,
		},
		[]byte{},
	)).Bytes()
	if !bytes.Equal(check, msg.Body().Hash()) {
		return nil, nil
	}

	// check size of data
	if len(dataBytes) < SizeUint64 {
		return nil, nil
	}

	// pass random bytes and get main data
	mustLen := encoding.BytesToUint64(dataBytes[:SizeUint64])
	allData := dataBytes[SizeUint64:]
	if mustLen > uint64(len(allData)) {
		return nil, nil
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.Decrypt(msg.Body().Sign())
	if sign == nil {
		return nil, nil
	}
	if !public.Verify(msg.Body().Hash(), sign) {
		return nil, nil
	}

	decMsg := &sMessage{
		FHead: sHeadMessage{
			FSender:  publicBytes,
			FSession: session,
			FSalt:    salt,
		},
		FBody: sBodyMessage{
			FData:  allData[:mustLen],
			FHash:  msg.Body().Hash(),
			FSign:  sign,
			FProof: msg.Body().Proof(),
		},
	}

	// export title from (title||data)
	title := decMsg.export()
	if title == nil {
		return nil, nil
	}

	// Return decrypted message with title.
	return decMsg, title
}
