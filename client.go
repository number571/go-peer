package gopeer

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/number571/gopeer/crypto"
	"github.com/number571/gopeer/encoding"
)

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv crypto.PrivKey) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		privateKey:  priv,
		hroutes:     make(map[string]func(*Client, *Package) []byte),
		mapping:     make(map[string]bool),
		connections: make(map[string]net.Conn),
		actions:     make(map[string]chan []byte),
		F2F: &friendToFriend{
			friends: make(map[string]crypto.PubKey),
		},
	}
}

// Create package: Head.Title = title, Body.Data = data.
func NewPackage(title, data []byte) *Package {
	return &Package{
		Head: HeadPackage{
			Title: []byte(title),
		},
		Body: BodyPackage{
			Data: data,
		},
	}
}

// Create route object with receiver.
func NewRoute(receiver crypto.PubKey) *Route {
	if receiver == nil {
		return nil
	}
	return &Route{
		receiver: receiver,
	}
}

// Append pseudo sender to route.
func (route *Route) Sender(psender crypto.PrivKey) *Route {
	route.psender = psender
	return route
}

// Append route-nodes to route.
func (route *Route) Routes(routes []crypto.PubKey) *Route {
	route.routes = routes
	return route
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
		id := encoding.Base64Encode(crypto.GenRand(settings.RAND_SIZE))
		client.setConnection(id, conn)
		go client.handleConn(id)
	}
	return nil
}

// Add function to mapping for route use.
func (client *Client) Handle(title []byte, handle func(*Client, *Package) []byte) *Client {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.hroutes[encoding.Base64Encode(title)] = handle
	return client
}

// Send package by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (client *Client) Send(pack *Package, route *Route) ([]byte, error) {
	var (
		err      error
		result   []byte
		hash     = string(route.receiver.Address())
		retryNum = settings.RETRY_NUM
	)

	client.setAction(hash)
	defer func() {
		client.delAction(hash)
	}()

repeat:
	routePack := client.RoutePackage(pack, route)
	if routePack == nil {
		return result, errors.New("psender is nil")
	}

	client.send(routePack)

	select {
	case result = <-client.actions[hash]:
	case <-time.After(time.Duration(settings.WAIT_TIME) * time.Second):
		if retryNum > 1 {
			retryNum -= 1
			goto repeat
		}
		err = errors.New("time is over")
	}

	return result, err
}

// Function wrap package in multiple route.
// Need use pseudo sender if route not null.
func (client *Client) RoutePackage(pack *Package, route *Route) *Package {
	var (
		rpack   = client.Encrypt(route.receiver, pack, settings.POWS_DIFF)
		psender = NewClient(route.psender)
	)
	if len(route.routes) != 0 && psender == nil {
		return nil
	}
	for _, pub := range route.routes {
		rpack = psender.Encrypt(
			pub,
			NewPackage(settings.ROUTE_MSG, serializePackage(rpack)),
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
	for _, addr := range addresses {
		if client.InConnections(addr) {
			client.getConnection(addr).Close()
		}
		client.delConnection(addr)
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

// Encrypt package with public key of receiver.
// The package can be decrypted only if private key is known.
func (client *Client) Encrypt(receiver crypto.PubKey, pack *Package, diff uint) *Package {
	var (
		rand = crypto.GenRand(uint(settings.RAND_SIZE))
		hash = crypto.HashSum(bytes.Join(
			[][]byte{
				rand,
				client.PubKey().Bytes(),
				receiver.Bytes(),
				[]byte(pack.Head.Title),
				pack.Body.Data,
			},
			[]byte{},
		))
		session = crypto.GenRand(uint(settings.SKEY_SIZE))
		cipher  = crypto.NewCipher(session)
	)
	return &Package{
		Head: HeadPackage{
			Rand:    cipher.Encrypt(rand),
			Title:   cipher.Encrypt(pack.Head.Title),
			Sender:  cipher.Encrypt(client.PubKey().Bytes()),
			Session: receiver.Encrypt(session),
		},
		Body: BodyPackage{
			Data: cipher.Encrypt(pack.Body.Data),
			Hash: hash,
			Sign: cipher.Encrypt(client.PrivKey().Sign(hash)),
			Npow: crypto.NewPuzzle(diff).Proof(hash),
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
	if !crypto.NewPuzzle(pow).Verify(hash, pack.Body.Npow) {
		return nil
	}

	session := client.PrivKey().Decrypt(pack.Head.Session)
	if session == nil {
		return nil
	}

	cipher := crypto.NewCipher(session)
	publicBytes := cipher.Decrypt(pack.Head.Sender)
	if publicBytes == nil {
		return nil
	}

	public := crypto.LoadPubKey(publicBytes)
	if public == nil {
		return nil
	}
	if public.Size() != settings.AKEY_SIZE {
		return nil
	}

	sign := cipher.Decrypt(pack.Body.Sign)
	if sign == nil {
		return nil
	}
	if !public.Verify(hash, sign) {
		return nil
	}

	titleBytes := cipher.Decrypt(pack.Head.Title)
	if titleBytes == nil {
		return nil
	}

	dataBytes := cipher.Decrypt(pack.Body.Data)
	if dataBytes == nil {
		return nil
	}

	rand := cipher.Decrypt(pack.Head.Rand)
	if rand == nil {
		return nil
	}

	check := crypto.HashSum(bytes.Join(
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

	return &Package{
		Head: HeadPackage{
			Rand:    rand,
			Title:   titleBytes,
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
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	return f2f.enabled
}

// Switch f2f mode to reverse.
func (f2f *friendToFriend) Switch() {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	f2f.enabled = !f2f.enabled
}

// Check the existence of a friend in the list by the public key.
func (f2f *friendToFriend) InList(pub crypto.PubKey) bool {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	_, ok := f2f.friends[string(pub.Address())]
	return ok
}

// Get a list of friends public keys.
func (f2f *friendToFriend) List() []crypto.PubKey {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	var list []crypto.PubKey
	for _, pub := range f2f.friends {
		list = append(list, pub)
	}
	return list
}

// Add public key to list of friends.
func (f2f *friendToFriend) Append(pubs ...crypto.PubKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	for _, pub := range pubs {
		f2f.friends[string(pub.Address())] = pub
	}
}

// Delete public key from list of friends.
func (f2f *friendToFriend) Remove(pubs ...crypto.PubKey) {
	f2f.mutex.Lock()
	defer f2f.mutex.Unlock()
	for _, pub := range pubs {
		delete(f2f.friends, string(pub.Address()))
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

	repeat:
		if pack == nil {
			continue
		}

		// size(sha256) = 32bytes
		if len(pack.Body.Hash) != 32 {
			continue
		}

		if client.inMapping(pack.Body.Hash) {
			continue
		}
		client.setMapping(pack.Body.Hash)

		if !crypto.NewPuzzle(settings.POWS_DIFF).Verify(pack.Body.Hash, pack.Body.Npow) {
			continue
		}
		client.send(pack)

		decPack := client.Decrypt(pack, settings.POWS_DIFF)
		if decPack == nil {
			continue
		}

		sender := crypto.LoadPubKey(decPack.Head.Sender)
		if client.F2F.State() && !client.F2F.InList(sender) {
			continue
		}

		if bytes.Equal(decPack.Head.Title, settings.ROUTE_MSG) {
			pack = deserializePackage(decPack.Body.Data)
			goto repeat
		}

		handleFunc(client, decPack)
	}
}

func handleFunc(client *Client, pack *Package) {
	fname := pack.Head.Title
	if bytes.HasPrefix(fname, settings.RET_BYTES) {
		client.response(
			crypto.LoadPubKey(pack.Head.Sender),
			pack.Body.Data,
		)
		return
	}
	client.send(client.Encrypt(
		crypto.LoadPubKey(pack.Head.Sender),
		NewPackage(
			bytes.Join([][]byte{
				settings.RET_BYTES,
				fname,
			}, []byte{}),
			client.getFunction(fname)(client, pack),
		),
		settings.POWS_DIFF,
	))
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
	return deserializePackage(message)
}

func (client *Client) send(pack *Package) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	bytesPack := bytes.Join(
		[][]byte{
			[]byte(serializePackage(pack)),
			[]byte(settings.END_BYTES),
		},
		[]byte{},
	)
	client.mapping[encoding.Base64Encode(pack.Body.Hash)] = true
	for _, cn := range client.connections {
		go cn.Write(bytesPack)
	}
}

func (client *Client) response(pub crypto.PubKey, data []byte) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	hash := string(pub.Address())
	if _, ok := client.actions[hash]; ok {
		client.actions[hash] <- data
	}
}

func (client *Client) getFunction(name []byte) func(*Client, *Package) []byte {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.hroutes[encoding.Base64Encode(name)]
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
	client.mapping[encoding.Base64Encode(hash)] = true
}

func (client *Client) inMapping(hash []byte) bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	_, ok := client.mapping[encoding.Base64Encode(hash)]
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

// Serialize with JSON format.
func serializePackage(pack *Package) []byte {
	jsonData, err := json.MarshalIndent(pack, "", "\t")
	if err != nil {
		return nil
	}
	return jsonData
}

// Deserialize with JSON format.
func deserializePackage(jsonData []byte) *Package {
	var pack = new(Package)
	err := json.Unmarshal(jsonData, pack)
	if err != nil {
		return nil
	}
	return pack
}
