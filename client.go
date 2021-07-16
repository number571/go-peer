package gopeer

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"math/big"
	"net"
	"time"
)

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv *rsa.PrivateKey) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		privateKey:  priv,
		functions:   make(map[string]func(*Client, *Package) []byte),
		mapping:     make(map[string]bool),
		connections: make(map[string]net.Conn),
		actions:     make(map[string]chan []byte),
		F2F: &friendToFriend{
			friends: make(map[string]*rsa.PublicKey),
		},
	}
}

// Create package: Head.Title = title, Body.Data = data.
func NewPackage(title string, data []byte) *Package {
	return &Package{
		Head: HeadPackage{
			Title: title,
		},
		Body: BodyPackage{
			Data: data,
		},
	}
}

// Turn on listener by address.
// Client handle function need be not null.
func (client *Client) RunNode(address string) error {
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
		if client.isMaxConnSize() {
			conn.Close()
			continue
		}
		id := Base64Encode(GenerateBytes(settings.RAND_SIZE))
		client.setConnection(id, conn)
		go client.handleConn(id)
	}
	return nil
}

// Add function to mapping for route use.
func (client *Client) Handle(title string, handle func(*Client, *Package) []byte) *Client {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.functions[title] = handle
	return client
}

// Send package by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (client *Client) Send(receiver *rsa.PublicKey, pack *Package, route []*rsa.PublicKey, ppsender *rsa.PrivateKey) ([]byte, error) {
	var (
		err      error
		result   []byte
		hash     = HashPublicKey(receiver)
		retryNum = settings.RETRY_NUM
	)

	client.setAction(hash)
	defer func() {
		client.delAction(hash)
	}()

tryAgain:
	routePack := client.RoutePackage(receiver, pack, route, ppsender)
	if routePack == nil {
		return result, errors.New("psender is nil")
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

// Function wrap package in multiple route.
// Need use pseudo sender if route not null.
func (client *Client) RoutePackage(receiver *rsa.PublicKey, pack *Package, route []*rsa.PublicKey, ppsender *rsa.PrivateKey) *Package {
	var (
		rpack   = client.Encrypt(receiver, pack, settings.POWS_DIFF)
		psender = NewClient(ppsender)
	)
	for _, pub := range route {
		if psender == nil {
			return nil
		}
		rpack = psender.Encrypt(
			pub,
			NewPackage(settings.ROUTE_MSG, SerializePackage(rpack)),
			settings.POWS_DIFF,
		)
	}
	return rpack
}

// Get list of connection addresses.
func (client *Client) Connections() []string {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	var list []string
	for addr := range client.connections {
		list = append(list, addr)
	}
	return list
}

// Check the existence of an address in the list of connections.
func (client *Client) InConnections(address string) bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	_, ok := client.connections[address]
	return ok
}

// Connect to node by address.
// Client handle function need be not null.
func (client *Client) Connect(addresses ...string) []error {
	var (
		listErrors []error = nil
	)
	for _, addr := range addresses {
		if client.isMaxConnSize() {
			return append(listErrors, errors.New("max conn"))
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			listErrors = append(listErrors, err)
			continue
		}
		client.setConnection(addr, conn)
		go client.handleConn(addr)
	}
	return listErrors
}

// Disconnect from node by address.
func (client *Client) Disconnect(addresses ...string) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	for _, addr := range addresses {
		if conn, ok := client.connections[addr]; ok {
			conn.Close()
		}
		delete(client.connections, addr)
	}
}

// Get public key from client object.
func (client *Client) PublicKey() *rsa.PublicKey {
	return &client.privateKey.PublicKey
}

// Get private key from client object.
func (client *Client) PrivateKey() *rsa.PrivateKey {
	return client.privateKey
}

// Encrypt package with public key of receiver.
// The package can be decrypted only if private key is known.
func (client *Client) Encrypt(receiver *rsa.PublicKey, pack *Package, pow uint) *Package {
	var (
		session = GenerateBytes(uint(settings.SKEY_SIZE))
		rand    = GenerateBytes(uint(settings.RAND_SIZE))
		hash    = HashSum(bytes.Join(
			[][]byte{
				rand,
				PublicKeyToBytes(client.PublicKey()),
				PublicKeyToBytes(receiver),
				[]byte(pack.Head.Title),
				pack.Body.Data,
			},
			[]byte{},
		))
		sign = Sign(client.PrivateKey(), hash)
	)
	return &Package{
		Head: HeadPackage{
			Rand:    EncryptAES(session, rand),
			Title:   Base64Encode(EncryptAES(session, []byte(pack.Head.Title))),
			Sender:  EncryptAES(session, PublicKeyToBytes(client.PublicKey())),
			Session: EncryptRSA(receiver, session),
		},
		Body: BodyPackage{
			Data: EncryptAES(session, pack.Body.Data),
			Hash: hash,
			Sign: EncryptAES(session, sign),
			Npow: ProofOfWork(hash, pow),
		},
	}
}

// Decrypt package with private key of receiver.
// No one else except the sender will be able to decrypt the package.
func (client *Client) Decrypt(pack *Package, pow uint) *Package {
	hash := pack.Body.Hash
	if hash == nil {
		return nil
	}
	if !ProofIsValid(hash, pow, pack.Body.Npow) {
		return nil
	}
	session := DecryptRSA(client.PrivateKey(), pack.Head.Session)
	if session == nil {
		return nil
	}
	publicBytes := DecryptAES(session, pack.Head.Sender)
	if publicBytes == nil {
		return nil
	}
	public := BytesToPublicKey(publicBytes)
	if public == nil {
		return nil
	}
	size := big.NewInt(1)
	size.Lsh(size, uint(settings.AKEY_SIZE-1))
	if public.N.Cmp(size) == -1 {
		return nil
	}
	sign := DecryptAES(session, pack.Body.Sign)
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
	dataBytes := DecryptAES(session, pack.Body.Data)
	if dataBytes == nil {
		return nil
	}
	rand := DecryptAES(session, pack.Head.Rand)
	if rand == nil {
		return nil
	}
	check := HashSum(bytes.Join(
		[][]byte{
			rand,
			publicBytes,
			PublicKeyToBytes(client.PublicKey()),
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
			Rand:    rand,
			Title:   string(titleBytes),
			Sender:  publicBytes,
			Session: session,
		},
		Body: BodyPackage{
			Data: dataBytes,
			Hash: hash,
			Sign: sign,
			Npow: pack.Body.Npow,
		},
	}
}

// Get current state of f2f mode.
func (f2f *friendToFriend) State() bool {
	return f2f.enabled
}

// Switch f2f mode to reverse.
func (f2f *friendToFriend) Switch() {
	f2f.enabled = !f2f.enabled
}

// Check the existence of a friend in the list by the public key.
func (f2f *friendToFriend) InList(pub *rsa.PublicKey) bool {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	_, ok := f2f.friends[HashPublicKey(pub)]
	return ok
}

// Get a list of friends public keys.
func (f2f *friendToFriend) List() []*rsa.PublicKey {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	var list []*rsa.PublicKey
	for _, pub := range f2f.friends {
		list = append(list, pub)
	}
	return list
}

// Add public key to list of friends.
func (f2f *friendToFriend) Append(pubs ...*rsa.PublicKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	for _, pub := range pubs {
		f2f.friends[HashPublicKey(pub)] = pub
	}
}

// Delete public key from list of friends.
func (f2f *friendToFriend) Remove(pubs ...*rsa.PublicKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	for _, pub := range pubs {
		delete(f2f.friends, HashPublicKey(pub))
	}
}

func (client *Client) handleConn(id string) {
	conn := client.getConnection(id)

	defer func() {
		conn.Close()
		client.delConnection(id)
	}()

	for {
		pack := readPackage(conn)

	checkAgain:
		if pack == nil {
			continue
		}

		// size(sha256) = 32 bytes
		if len(pack.Body.Hash) != 32 {
			continue
		}

		if client.inMapping(pack.Body.Hash) {
			continue
		}
		client.setMapping(pack.Body.Hash)

		if !ProofIsValid(pack.Body.Hash, settings.POWS_DIFF, pack.Body.Npow) {
			continue
		}
		client.send(pack)

		decPack := client.Decrypt(pack, settings.POWS_DIFF)
		if decPack == nil {
			continue
		}

		if client.F2F.State() && !client.F2F.InList(BytesToPublicKey(decPack.Head.Sender)) {
			continue
		}

		if decPack.Head.Title == settings.ROUTE_MSG {
			pack = DeserializePackage(decPack.Body.Data)
			goto checkAgain
		}

		handleFunc(client, decPack)
	}
}

func handleFunc(client *Client, pack *Package) {
	for nm, fn := range client.functions {
		switch pack.Head.Title {
		case nm:
			client.send(client.Encrypt(
				BytesToPublicKey(pack.Head.Sender),
				NewPackage("_"+nm, fn(client, pack)),
				settings.POWS_DIFF,
			))
			return
		case "_" + nm:
			client.response(
				BytesToPublicKey(pack.Head.Sender),
				pack.Body.Data,
			)
			return
		}
	}
}

func readPackage(conn net.Conn) *Package {
	var (
		message []byte
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
		message = bytes.Join(
			[][]byte{
				message,
				buffer[:length],
			},
			[]byte{},
		)
		if bytes.Contains(message, []byte(settings.END_BYTES)) {
			message = bytes.Split(message, []byte(settings.END_BYTES))[0]
			break
		}
	}
	return DeserializePackage(message)
}

func (client *Client) send(pack *Package) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	bytesPack := bytes.Join(
		[][]byte{
			[]byte(SerializePackage(pack)),
			[]byte(settings.END_BYTES),
		},
		[]byte{},
	)
	client.mapping[Base64Encode(pack.Body.Hash)] = true
	for _, cn := range client.connections {
		go cn.Write(bytesPack)
	}
}

func (client *Client) response(pub *rsa.PublicKey, data []byte) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	hash := HashPublicKey(pub)
	if _, ok := client.actions[hash]; ok {
		client.actions[hash] <- data
	}
}

func (client *Client) setAction(hash string) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.actions[hash] = make(chan []byte)
}

func (client *Client) delAction(hash string) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	delete(client.actions, hash)
}

func (client *Client) setMapping(hash []byte) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if uint(len(client.mapping)) > settings.MAPP_SIZE {
		client.mapping = make(map[string]bool)
	}
	client.mapping[Base64Encode(hash)] = true
}

func (client *Client) inMapping(hash []byte) bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	_, ok := client.mapping[Base64Encode(hash)]
	return ok
}

func (client *Client) isMaxConnSize() bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return uint(len(client.connections)) > settings.CONN_SIZE
}

func (client *Client) setConnection(id string, conn net.Conn) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.connections[id] = conn
}

func (client *Client) getConnection(id string) net.Conn {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.connections[id]
}

func (client *Client) delConnection(id string) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	delete(client.connections, id)
}
