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

func NewClient(priv *rsa.PrivateKey) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		mutex:       new(sync.Mutex),
		privateKey: priv,
		mapping:     make(map[string]bool),
		connections: make(map[net.Conn]string),
		actions:    make(map[string]chan string),
		f2f: friendToFriend{
			friends: make(map[string]*rsa.PublicKey),
		},
	}
}

func (client *Client) Send(receiver *rsa.PublicKey, pack *Package) (string, error) {
	var (
		err    error
		result string
		hash   = HashPublic(receiver)
	)

	client.actions[hash] = make(chan string)
	defer delete(client.actions, hash)

	client.send(receiver, pack)

	select {
	case result = <-client.actions[hash]:
	case <-time.After(time.Duration(settings.WAIT_TIME) * time.Second):
		err = errors.New("time is over")
	}

	return result, err
}

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

func (client *Client) Public() *rsa.PublicKey {
	return &client.privateKey.PublicKey
}

func (client *Client) Private() *rsa.PrivateKey {
	return client.privateKey
}

func (client *Client) StringPublic() string {
	return StringPublic(&client.privateKey.PublicKey)
}

func (client *Client) StringPrivate() string {
	return StringPrivate(client.privateKey)
}

func (client *Client) HashPublic() string {
	return HashPublic(&client.privateKey.PublicKey)
}

func (client *Client) F2F() bool {
	return client.f2f.enabled
}

func (client *Client) EnableF2F() {
	client.f2f.enabled = true
}

func (client *Client) DisableF2F() {
	client.f2f.enabled = false
}

func (client *Client) InF2F(pub *rsa.PublicKey) bool {
	if _, ok := client.f2f.friends[HashPublic(pub)]; ok {
		return true
	}
	return false
}

func (client *Client) ListF2F() []rsa.PublicKey {
	var list []rsa.PublicKey
	for _, pub := range client.f2f.friends {
		list = append(list, *pub)
	}
	return list
}

func (client *Client) AppendF2F(pub *rsa.PublicKey) {
	client.f2f.friends[HashPublic(pub)] = pub
}

func (client *Client) RemoveF2F(pub *rsa.PublicKey) {
	delete(client.f2f.friends, HashPublic(pub))
}

func (client *Client) send(receiver *rsa.PublicKey, pack *Package) {
	encPack := client.encrypt(receiver, pack)
	bytesPack := EncodePackage(encPack)
	client.mapping[encPack.Body.Hash] = true
	for cn := range client.connections {
		go cn.Write(bytes.Join(
			[][]byte{
				[]byte(bytesPack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		))
	}
}

func (client *Client) redirect(pack *Package, sender net.Conn) {
	encPack := EncodePackage(pack)
	for cn := range client.connections {
		if cn == sender {
			continue
		}
		go cn.Write(bytes.Join(
			[][]byte{
				[]byte(encPack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		))
	}
}

func (client *Client) response(pub *rsa.PublicKey, data string) {
	client.mutex.Lock()
	hash := HashPublic(pub)
	if _, ok := client.actions[hash]; ok {
		client.actions[hash] <- data
	}
	client.mutex.Unlock()
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
			Npow: ProofOfWork(hash, settings.POWS_DIFF),
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
	if !ProofIsValid(hash, pack.Body.Npow) {
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
			Npow: pack.Body.Npow,
		},
	}
}
