package gopeer

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"math/big"
	"net"
	"sync"
	"time"
)


// CREATE
func NewClient(priv *rsa.PrivateKey) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		mutex:       new(sync.Mutex),
		mapping:     make(map[string]bool),
		connections: make(map[net.Conn]string),
		publicKey:   &priv.PublicKey,
		privateKey:  priv,
		actions:     make(map[string]chan bool),
		F2F: FriendToFriend{
			friends: make(map[string]*rsa.PublicKey),
		},
	}
}


// SEND / REQUEST / RESPONSE
func (client *Client) Send(receiver *rsa.PublicKey, pack *Package) {
	client.mapping[pack.Body.Hash] = true
	encPack := EncodePackage(client.encrypt(receiver, pack))
	for cn := range client.connections {
		cn.Write(bytes.Join(
			[][]byte{
				[]byte(encPack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		))
	}
}

func (client *Client) Request(receiver *rsa.PublicKey, pack *Package) error {
	hash := HashPublic(receiver)

	client.actions[hash] = make(chan bool)
	defer delete(client.actions, hash)

	client.Send(receiver, pack)

	select {
	case <-client.actions[hash]:
	case <-time.After(time.Duration(settings.WAIT_TIME) * time.Second):
		return errors.New("time is over")
	}

	return nil
}

func (client *Client) Response(pub *rsa.PublicKey) {
	hash := HashPublic(pub)
	if _, ok := client.actions[hash]; ok {
		client.actions[hash] <- true
	}
}


// CONNECT / DISCONNECT
func (client *Client) Connect(address string, handle func(*Client, *Package)) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	client.connections[conn] = address
	go handleConn(conn, client, handle)
	return nil
}

func (client *Client) Disconnect(address string) {
	for conn, addr := range client.connections {
		if addr == address {
			delete(client.connections, conn)
			conn.Close()
		}
	}
}


// PUBLIC / PRIVATE
func (client *Client) Public() *rsa.PublicKey {
	return client.publicKey
}

func (client *Client) Private() *rsa.PrivateKey {
	return client.privateKey
}

func (client *Client) StringPublic() string {
	return StringPublic(client.publicKey)
}

func (client *Client) StringPrivate() string {
	return StringPrivate(client.privateKey)
}

func (client *Client) HashPublic() string {
	return HashPublic(client.publicKey)
}


// F2F
func (client *Client) AppendToF2F(pub *rsa.PublicKey) {
	client.F2F.friends[HashPublic(pub)] = pub
}

func (client *Client) RemoveFromF2F(pub *rsa.PublicKey) {
	delete(client.F2F.friends, HashPublic(pub))
}

func (client *Client) InF2F(pub *rsa.PublicKey) bool {
	if _, ok := client.F2F.friends[HashPublic(pub)]; ok {
		return true
	}
	return false
}


// LOCAL DATA
func (client *Client) redirect(pack *Package, sender net.Conn) {
	encPack := EncodePackage(pack)
	for cn := range client.connections {
		if cn == sender {
			continue
		}
		cn.Write(bytes.Join(
			[][]byte{
				[]byte(encPack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		))
	}
}

func (client *Client) encrypt(receiver *rsa.PublicKey, pack *Package) *Package {
	var (
		session = GenerateBytes(uint(settings.SKEY_SIZE))
		rand    = GenerateBytes(uint(settings.RAND_SIZE))
		hash    = HashSum(bytes.Join(
			[][]byte{
				rand,
				Base64Decode(client.StringPublic()),
				Base64Decode(StringPublic(receiver)),
				[]byte(pack.Head.Title),
				[]byte(pack.Body.Data),
			},
			[]byte{},
		))
		sign = Sign(client.privateKey, hash)
	)
	return &Package{
		Head: HeadPackage{
			Rand:    Base64Encode(EncryptAES(session, rand)),
			Title:   Base64Encode(EncryptAES(session, []byte(pack.Head.Title))),
			Sender:  Base64Encode(EncryptAES(session, Base64Decode(client.StringPublic()))),
			Session: Base64Encode(EncryptRSA(receiver, session)),
		},
		Body: BodyPackage{
			Data: Base64Encode(EncryptAES(session, []byte(pack.Body.Data))),
			Hash: Base64Encode(hash),
			Sign: Base64Encode(sign),
		},
	}
}

func (client *Client) decrypt(pack *Package) *Package {
	session := DecryptRSA(client.privateKey, Base64Decode(pack.Head.Session))
	if session == nil {
		return nil
	}
	publicBytes := DecryptAES(session, Base64Decode(pack.Head.Sender))
	if publicBytes == nil {
		return nil
	}
	public := ParsePublic(Base64Encode(publicBytes))
	if public == nil {
		return nil
	}
	size := big.NewInt(1)
	size.Lsh(size, uint(settings.AKEY_SIZE-1))
	if public.N.Cmp(size) == -1 {
		return nil
	}
	titleBytes := DecryptAES(session, Base64Decode(pack.Head.Title))
	if titleBytes == nil {
		return nil
	}
	dataBytes := DecryptAES(session, Base64Decode(pack.Body.Data))
	if dataBytes == nil {
		return nil
	}
	rand := DecryptAES(session, Base64Decode(pack.Head.Rand))
	hash := HashSum(bytes.Join(
		[][]byte{
			rand,
			publicBytes,
			Base64Decode(client.StringPublic()),
			titleBytes,
			dataBytes,
		},
		[]byte{},
	))
	if Base64Encode(hash) != pack.Body.Hash {
		return nil
	}
	err := Verify(public, hash, Base64Decode(pack.Body.Sign))
	if err != nil {
		return nil
	}
	return &Package{
		Head: HeadPackage{
			Rand:    Base64Encode(rand),
			Title:   string(titleBytes),
			Sender:  Base64Encode(publicBytes),
			Session: Base64Encode(session),
		},
		Body: BodyPackage{
			Data: string(dataBytes),
			Hash: pack.Body.Hash,
			Sign: pack.Body.Sign,
		},
	}
}
