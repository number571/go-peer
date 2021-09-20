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
		rand = crypto.Rand(gopeer.Get("RAND_SIZE").(uint))
		hash = crypto.SumHash(bytes.Join(
			[][]byte{
				rand,
				client.PubKey().Bytes(),
				receiver.Bytes(),
				[]byte(msg.Head.Title),
				msg.Body.Data,
			},
			[]byte{},
		))
		session = crypto.Rand(gopeer.Get("SKEY_SIZE").(uint))
		cipher  = crypto.NewCipher(session)
	)
	return &Message{
		Head: HeadMessage{
			diff:    msg.Head.diff,
			Rand:    cipher.Encrypt(rand),
			Title:   cipher.Encrypt(msg.Head.Title),
			Sender:  cipher.Encrypt(client.PubKey().Bytes()),
			Session: receiver.Encrypt(session),
		},
		Body: BodyMessage{
			Data: cipher.Encrypt(msg.Body.Data),
			Hash: hash,
			Sign: cipher.Encrypt(client.PrivKey().Sign(hash)),
			Npow: crypto.NewPuzzle(msg.Head.diff).Proof(hash),
		},
	}
}

// Decrypt message with private key of receiver.
// No one else except the sender will be able to decrypt the message.
func (client *Client) Decrypt(msg *Message) *Message {
	hash := msg.Body.Hash

	if !crypto.NewPuzzle(msg.Head.diff).Verify(hash, msg.Body.Npow) {
		return nil
	}

	session := client.PrivKey().Decrypt(msg.Head.Session)
	if session == nil {
		return nil
	}

	cipher := crypto.NewCipher(session)
	publicBytes := cipher.Decrypt(msg.Head.Sender)
	if publicBytes == nil {
		return nil
	}

	public := crypto.LoadPubKey(publicBytes)
	if public == nil {
		return nil
	}
	if public.Size() != gopeer.Get("AKEY_SIZE").(uint) {
		return nil
	}

	sign := cipher.Decrypt(msg.Body.Sign)
	if sign == nil {
		return nil
	}
	if !public.Verify(hash, sign) {
		return nil
	}

	titleBytes := cipher.Decrypt(msg.Head.Title)
	if titleBytes == nil {
		return nil
	}

	dataBytes := cipher.Decrypt(msg.Body.Data)
	if dataBytes == nil {
		return nil
	}

	rand := cipher.Decrypt(msg.Head.Rand)
	if rand == nil {
		return nil
	}

	check := crypto.SumHash(bytes.Join(
		[][]byte{
			rand,
			publicBytes,
			client.PubKey().Bytes(),
			titleBytes,
			dataBytes,
		},
		[]byte{},
	))
	if !bytes.Equal(check, hash) {
		return nil
	}

	return &Message{
		Head: HeadMessage{
			diff:    msg.Head.diff,
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
	if len(route.routes) != 0 && psender == nil {
		return nil
	}
	diff := uint(msg.Head.diff)
	pack := rmsg.Serialize()
	for _, pub := range route.routes {
		rmsg = psender.Encrypt(
			pub,
			NewMessage(
				gopeer.Get("ROUTE_MSG").([]byte),
				pack.Bytes(),
			).WithDiff(diff),
		)
	}
	return rmsg
}
