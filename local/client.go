package local

import (
	"bytes"

	"github.com/number571/gopeer"
	"github.com/number571/gopeer/crypto"
)

// Basic structure describing the user.
type Client struct {
	privateKey crypto.PrivKey
}

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv crypto.PrivKey) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		privateKey: priv,
	}
}

// Get public key from client object.
func (client *Client) PubKey() crypto.PubKey {
	return client.privateKey.PubKey()
}

// Get private key from client object.
func (client *Client) PrivKey() crypto.PrivKey {
	return client.privateKey
}

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *Client) Encrypt(receiver crypto.PubKey, msg *Message) *Message {
	var (
		rand = crypto.RandBytes(gopeer.Get("SALT_SIZE").(uint))
		hash = crypto.NewSHA256(bytes.Join(
			[][]byte{
				rand,
				client.PubKey().Bytes(),
				receiver.Bytes(),
				[]byte(msg.Head.Title),
				msg.Body.Data,
			},
			[]byte{},
		)).Bytes()
		session = crypto.RandBytes(gopeer.Get("SKEY_SIZE").(uint))
		cipher  = crypto.NewCipher(session)
	)

	return &Message{
		Head: HeadMessage{
			Diff:    msg.Head.Diff,
			Rand:    cipher.Encrypt(rand),
			Title:   cipher.Encrypt(msg.Head.Title),
			Sender:  cipher.Encrypt(client.PubKey().Bytes()),
			Session: receiver.Encrypt(session),
		},
		Body: BodyMessage{
			Data: cipher.Encrypt(msg.Body.Data),
			Hash: hash,
			Sign: cipher.Encrypt(client.PrivKey().Sign(hash)),
			Npow: crypto.NewPuzzle(msg.Head.Diff).Proof(hash),
		},
	}
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *Client) Decrypt(msg *Message) *Message {
	hash := msg.Body.Hash

	// Proof of work. Prevent spam.
	puzzle := crypto.NewPuzzle(msg.Head.Diff)
	if !puzzle.Verify(hash, msg.Body.Npow) {
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
	if public.Size() != gopeer.Get("AKEY_SIZE").(uint) {
		return nil
	}

	// Decrypt sign of message and verify this
	// by public key of sender and hash of message.
	sign := cipher.Decrypt(msg.Body.Sign)
	if sign == nil {
		return nil
	}
	if !public.Verify(hash, sign) {
		return nil
	}

	// Decrypt title of message by session key.
	titleBytes := cipher.Decrypt(msg.Head.Title)
	if titleBytes == nil {
		return nil
	}

	// Decrypt main data of message by session key.
	dataBytes := cipher.Decrypt(msg.Body.Data)
	if dataBytes == nil {
		return nil
	}

	// Decrypt random string by session key.
	rand := cipher.Decrypt(msg.Head.Rand)
	if rand == nil {
		return nil
	}

	// Check received hash and generated hash.
	check := crypto.NewSHA256(bytes.Join(
		[][]byte{
			rand,
			publicBytes,
			client.PubKey().Bytes(),
			titleBytes,
			dataBytes,
		},
		[]byte{},
	)).Bytes()
	if !bytes.Equal(check, hash) {
		return nil
	}

	// Return decrypted message.
	return &Message{
		Head: HeadMessage{
			Diff:    msg.Head.Diff,
			Title:   titleBytes,
			Rand:    rand,
			Sender:  publicBytes,
			Session: session,
		},
		Body: BodyMessage{
			Data: dataBytes,
			Hash: hash,
			Sign: sign,
			Npow: msg.Body.Npow,
		},
	}
}

// Function wrap message in multiple route.
// Need use pseudo sender if route not null.
func (client *Client) RouteMessage(msg *Message, route *Route) *Message {
	var (
		rmsg    = client.Encrypt(route.receiver, msg)
		psender = NewClient(route.psender)
	)
	if psender == nil && len(route.routes) != 0 {
		return nil
	}
	diff := uint(msg.Head.Diff)
	for _, pub := range route.routes {
		pack := rmsg.Serialize()
		rmsg = psender.Encrypt(
			pub,
			NewMessage(
				gopeer.Get("ROUTE_MSG").([]byte),
				pack.Bytes(),
				diff,
			),
		)
	}
	return rmsg
}
