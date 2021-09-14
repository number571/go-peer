package network

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/number571/gopeer"
	"github.com/number571/gopeer/crypto"
	"github.com/number571/gopeer/encoding"
)

// Basic structure describing the user.
// Stores the private key and list of friends.
type Client struct {
	mutex       sync.Mutex
	privateKey  crypto.PrivKey
	hroutes     map[string]func(*Client, *Message) []byte
	mapping     map[string]bool
	connections map[string]net.Conn
	actions     map[string]chan []byte
	F2F         *friendToFriend
}

type friendToFriend struct {
	mutex   sync.Mutex
	enabled bool
	friends map[string]crypto.PubKey
}

// Create client by private key as identification.
// Handle function is used when the network exists. Can be null.
func NewClient(priv crypto.PrivKey) *Client {
	if priv == nil {
		return nil
	}
	return &Client{
		privateKey:  priv,
		hroutes:     make(map[string]func(*Client, *Message) []byte),
		mapping:     make(map[string]bool),
		connections: make(map[string]net.Conn),
		actions:     make(map[string]chan []byte),
		F2F: &friendToFriend{
			friends: make(map[string]crypto.PubKey),
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
		id := encoding.Base64Encode(crypto.GenRand(gopeer.Get("RAND_SIZE").(uint)))
		client.setConnection(id, conn)
		go client.handleConn(id)
	}
	return nil
}

// Add function to mapping for route use.
func (client *Client) Handle(title []byte, handle func(*Client, *Message) []byte) *Client {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.hroutes[encoding.Base64Encode(title)] = handle
	return client
}

// Send message by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (client *Client) Send(msg *Message, route *Route) ([]byte, error) {
	var (
		err      error
		result   []byte
		hash     = string(route.receiver.Address())
		retryNum = gopeer.Get("RETRY_NUM").(uint)
	)

	client.setAction(hash)
	defer func() {
		client.delAction(hash)
	}()

REPEAT:
	routeMsg := client.RouteMessage(msg, route)
	if routeMsg == nil {
		return result, errors.New("psender is nil")
	}

	client.send(routeMsg)

	select {
	case result = <-client.actions[hash]:
	case <-time.After(time.Duration(gopeer.Get("WAIT_TIME").(uint)) * time.Second):
		if retryNum > 1 {
			retryNum -= 1
			goto REPEAT
		}
		err = errors.New("time is over")
	}

	return result, err
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
	diff := uint(rmsg.Head.Diff)
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

// Encrypt message with public key of receiver.
// The message can be decrypted only if private key is known.
func (client *Client) Encrypt(receiver crypto.PubKey, msg *Message) *Message {
	var (
		rand = crypto.GenRand(gopeer.Get("RAND_SIZE").(uint))
		hash = crypto.HashSum(bytes.Join(
			[][]byte{
				rand,
				client.PubKey().Bytes(),
				receiver.Bytes(),
				[]byte(msg.Head.Title),
				msg.Body.Data,
			},
			[]byte{},
		))
		session = crypto.GenRand(gopeer.Get("SKEY_SIZE").(uint))
		cipher  = crypto.NewCipher(session)
	)
	return &Message{
		Head: HeadMessage{
			Rand:    cipher.Encrypt(rand),
			Diff:    msg.Head.Diff,
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
	if hash == nil {
		return nil
	}
	if !crypto.NewPuzzle(msg.Head.Diff).Verify(hash, msg.Body.Npow) {
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

	return &Message{
		Head: HeadMessage{
			Title:   titleBytes,
			Diff:    msg.Head.Diff,
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
		msg := readMessage(conn)

	REPEAT:
		if msg == nil {
			continue
		}

		// size(sha256) = 32bytes
		if len(msg.Body.Hash) != 32 {
			continue
		}

		if client.inMapping(msg.Body.Hash) {
			continue
		}
		client.setMapping(msg.Body.Hash)

		puzzle := crypto.NewPuzzle(uint8(gopeer.Get("POWS_DIFF").(uint)))
		if !puzzle.Verify(msg.Body.Hash, msg.Body.Npow) {
			continue
		}
		client.send(msg)

		decMsg := client.Decrypt(msg)
		if decMsg == nil {
			continue
		}

		sender := crypto.LoadPubKey(decMsg.Head.Sender)
		if client.F2F.State() && !client.F2F.InList(sender) {
			continue
		}

		if bytes.Equal(decMsg.Head.Title, gopeer.Get("ROUTE_MSG").([]byte)) {
			msg = Package(decMsg.Body.Data).Deserialize()
			goto REPEAT
		}

		client.handleFunc(decMsg)
	}
}

func (client *Client) handleFunc(msg *Message) {
	fname := msg.Head.Title
	if bytes.HasPrefix(fname, gopeer.Get("RET_BYTES").([]byte)) {
		client.response(
			crypto.LoadPubKey(msg.Head.Sender),
			msg.Body.Data,
		)
		return
	}
	diff := uint(msg.Head.Diff)
	client.send(client.Encrypt(
		crypto.LoadPubKey(msg.Head.Sender),
		NewMessage(
			bytes.Join([][]byte{
				gopeer.Get("RET_BYTES").([]byte),
				fname,
			}, []byte{}),
			client.getFunction(fname)(client, msg),
		).WithDiff(diff),
	))
}

func (client *Client) send(msg *Message) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	pack := msg.Serialize()
	bytesMsg := bytes.Join(
		[][]byte{
			pack.Size(),
			pack.Bytes(),
		},
		[]byte{},
	)
	client.mapping[encoding.Base64Encode(msg.Body.Hash)] = true
	for _, cn := range client.connections {
		go cn.Write(bytesMsg)
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

func (client *Client) getFunction(name []byte) func(*Client, *Message) []byte {
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
	if uint(len(client.mapping)) > gopeer.Get("MAPP_SIZE").(uint) {
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
	return uint(len(client.connections)) > gopeer.Get("CONN_SIZE").(uint)
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

func readMessage(conn net.Conn) *Message {
	const (
		UINT64_SIZE = 8 // bytes
	)
	var (
		pack   []byte
		size   = uint(0)
		buflen = make([]byte, UINT64_SIZE)
		buffer = make([]byte, gopeer.Get("BUFF_SIZE").(uint))
	)

	length, err := conn.Read(buflen)
	if err != nil {
		return nil
	}
	if length != UINT64_SIZE {
		return nil
	}

	mustLen := uint(encoding.BytesToUint64(buflen))
	if mustLen > gopeer.Get("PACK_SIZE").(uint) {
		return nil
	}

	for {
		length, err = conn.Read(buffer)
		if err != nil {
			return nil
		}

		size += uint(length)
		if size > mustLen {
			return nil
		}

		pack = bytes.Join(
			[][]byte{
				pack,
				buffer[:length],
			},
			[]byte{},
		)

		if size == mustLen {
			break
		}
	}

	return Package(pack).Deserialize()
}
