package gopeer

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"
)

func NewClient(priv *rsa.PrivateKey, handle func(*Client, *Package)) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		handle:		 handle,
		mutex:       new(sync.Mutex),
		privateKey:  priv,
		mapping:     make(map[string]bool),
		connections: make(map[net.Conn]string),
		actions:     make(map[string]chan string),
		f2f: friendToFriend{
			friends: make(map[string]*rsa.PublicKey),
		},
	}
}

func NewPackage(title, data string) *Package {
	return &Package{
		Head: HeadPackage {
			Title: title,
		},
		Body: BodyPackage {
			Data: data,
		},
	}
}

func (client *Client) RunNode(address string) error {
	if client.handle == nil {
		return errors.New("handle nil")
	}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		client.mutex.Lock()
		if uint(len(client.connections)) > settings.CONN_SIZE {
			client.mutex.Unlock()
			conn.Close()
			continue
		}
		client.connections[conn] = "client"
		client.mutex.Unlock()
		go client.handleConn(conn, client.handle)
	}
	return nil
}

func (client *Client) Handle(title string, pack *Package, handle func(*Client, *Package) string) {
	switch pack.Head.Title {
	case title:
		client.send(client.Encrypt(
			ParsePublic(pack.Head.Sender), 
			NewPackage("_" + title, handle(client, pack)),
		))
	case "_" + title:
		client.response(
			ParsePublic(pack.Head.Sender),
			pack.Body.Data,
		)
	}
}

func (client *Client) Send(receiver *rsa.PublicKey, pack *Package, route []*rsa.PublicKey, pseudoSender *Client) (string, error) {
	var (
		err      error
		result   string
		hash     = HashPublic(receiver)
		retryNum = settings.RETRY_NUM
	)

	client.actions[hash] = make(chan string)
	defer func() {
		client.mutex.Lock()
		delete(client.actions, hash)
		client.mutex.Unlock()
	}()

tryAgain:

	routePack := client.Encrypt(receiver, pack)

	if pseudoSender == nil {
		pseudoSender = client
	}

	for _, pub := range route {
		routePack = pseudoSender.Encrypt(pub, 
			NewPackage(settings.ROUTE_MSG, SerializePackage(routePack)))
	}

	client.send(routePack)

	select {
	case result = <-client.actions[hash]:
	case <-time.After(time.Duration(settings.WAIT_TIME) * time.Second):
		if retryNum > 1 {
			retryNum -= 1
			goto tryAgain
		}
		err = errors.New("time is over")
	}

	return result, err
}

func (client *Client) Connect(address string) error {
	if client.handle == nil {
		return errors.New("handle nil")
	}
	if uint(len(client.connections)) > settings.CONN_SIZE {
		return errors.New("max conn")
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	client.mutex.Lock()
	client.connections[conn] = address
	client.mutex.Unlock()
	go client.handleConn(conn, client.handle)
	return nil
}

func (client *Client) Disconnect(address string) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
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
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if _, ok := client.f2f.friends[HashPublic(pub)]; ok {
		return true
	}
	return false
}

func (client *Client) ListF2F() []rsa.PublicKey {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	var list []rsa.PublicKey
	for _, pub := range client.f2f.friends {
		list = append(list, *pub)
	}
	return list
}

func (client *Client) AppendF2F(pub *rsa.PublicKey) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.f2f.friends[HashPublic(pub)] = pub
}

func (client *Client) RemoveF2F(pub *rsa.PublicKey) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	delete(client.f2f.friends, HashPublic(pub))
}

func (client *Client) Encrypt(receiver *rsa.PublicKey, pack *Package) *Package {
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
			Sign: Base64Encode(EncryptAES(session, sign)),
			Npow: ProofOfWork(hash, settings.POWS_DIFF),
		},
	}
}

func (client *Client) Decrypt(pack *Package) *Package {
	hash := Base64Decode(pack.Body.Hash)
	if hash == nil {
		return nil
	}
	if !ProofIsValid(hash, settings.POWS_DIFF, pack.Body.Npow) {
		return nil
	}
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
	sign := DecryptAES(session, Base64Decode(pack.Body.Sign))
	if sign == nil {
		return nil
	}
	err := Verify(public, hash, sign)
	if err != nil {
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
	check := HashSum(bytes.Join(
		[][]byte{
			rand,
			publicBytes,
			Base64Decode(client.StringPublic()),
			titleBytes,
			dataBytes,
		},
		[]byte{},
	))
	if !bytes.Equal(check, hash) {
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
			Sign: Base64Encode(sign),
			Npow: pack.Body.Npow,
		},
	}
}

func (client *Client) handleConn(conn net.Conn, handle func(*Client, *Package)) {
	defer func() {
		conn.Close()
		client.mutex.Lock()
		delete(client.connections, conn)
		client.mutex.Unlock()
	}()
	for {
		pack := readPackage(conn)
		isRoute := false

checkAgain:

		if pack == nil {
			continue
		}

		client.mutex.Lock()
		if _, ok := client.mapping[pack.Body.Hash]; ok {
			client.mutex.Unlock()
			continue
		}
		if uint(len(client.mapping)) > settings.MAPP_SIZE {
			client.mapping = make(map[string]bool)
		}
		client.mapping[pack.Body.Hash] = true
		client.mutex.Unlock()

		if !ProofIsValid(Base64Decode(pack.Body.Hash), settings.POWS_DIFF, pack.Body.Npow) {
			continue
		}

		if isRoute {
			client.send(pack)
		} else {
			client.redirect(pack, conn)
		}
		
		decPack := client.Decrypt(pack)

		if decPack == nil {
			continue
		}
		
		if client.f2f.enabled && !client.InF2F(ParsePublic(decPack.Head.Sender)) {
			continue
		}

		if decPack.Head.Title == settings.ROUTE_MSG {
			pack = DeserializePackage(decPack.Body.Data)
			isRoute = true
			goto checkAgain
		}

		handle(client, decPack)
	}
}

func readPackage(conn net.Conn) *Package {
	var (
		message string
		size    = uint(0)
		buffer  = make([]byte, settings.BUFF_SIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			return nil
		}
		size += uint(length)
		if size > settings.PACK_SIZE {
			return nil
		}
		message += string(buffer[:length])
		if strings.Contains(message, settings.END_BYTES) {
			message = strings.Split(message, settings.END_BYTES)[0]
			break
		}
	}
	return DeserializePackage(message)
}

func (client *Client) send(pack *Package) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	bytesPack := SerializePackage(pack)
	client.mapping[pack.Body.Hash] = true
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
	client.mutex.Lock()
	defer client.mutex.Unlock()
	encPack := SerializePackage(pack)
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
	defer client.mutex.Unlock()
	hash := HashPublic(pub)
	if _, ok := client.actions[hash]; ok {
		client.actions[hash] <- data
	}
}
