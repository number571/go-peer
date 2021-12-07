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
	"github.com/number571/gopeer/local"
)

// Basic structure for network use.
type Node struct {
	client      *local.Client
	f2f         *friendToFriend
	mutex       sync.Mutex
	hroutes     map[string]func(*local.Client, *local.Message) []byte
	mapping     map[string]bool
	connections map[string]net.Conn
	actions     map[string]chan []byte
}

// Create client by private key as identification.
func NewNode(client *local.Client) *Node {
	if client == nil {
		return nil
	}

	return &Node{
		client:      client,
		hroutes:     make(map[string]func(*local.Client, *local.Message) []byte),
		mapping:     make(map[string]bool),
		connections: make(map[string]net.Conn),
		actions:     make(map[string]chan []byte),
		f2f: &friendToFriend{
			friends: make(map[string]crypto.PubKey),
		},
	}
}

// Return client structure.
func (node *Node) Client() *local.Client {
	return node.client
}

// Return f2f structure.
func (node *Node) F2F() *friendToFriend {
	return node.f2f
}

// Turn on listener by address.
// Client handle function need be not null.
func (node *Node) Listen(address string) error {
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

		if node.hasMaxConnSize() {
			conn.Close()
			continue
		}

		id := crypto.RandString(gopeer.Get("SALT_SIZE").(uint))

		node.setConnection(id, conn)
		go node.handleConn(id)
	}

	return nil
}

// Add function to mapping for route use.
func (node *Node) Handle(title []byte, handle func(*local.Client, *local.Message) []byte) *Node {
	node.setFunction(title, handle)
	return node
}

// Send message by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (node *Node) Send(msg *local.Message, route *local.Route) ([]byte, error) {
	var (
		result []byte
		err    error
	)

	var (
		nonce    = crypto.RandBytes(gopeer.Get("SALT_SIZE").(uint))
		waitTime = time.Duration(gopeer.Get("WAIT_TIME").(uint))
		retryNum = gopeer.Get("RETRY_NUM").(uint)
		counter  = uint(0)
	)

	copyMsg := *msg
	copyMsg.Head.Title = bytes.Join(
		[][]byte{
			nonce,
			msg.Head.Title,
		},
		[]byte{},
	)

	node.setAction(nonce)
	defer node.delAction(nonce)

LOOP:
	for counter = 0; counter < retryNum; counter++ {
		routeMsg := node.client.RouteMessage(&copyMsg, route)
		if routeMsg == nil {
			return nil, errors.New("psender is nil and routes not nil")
		}

		node.send(routeMsg)

		select {
		case result = <-node.getAction(nonce):
			break LOOP
		case <-time.After(waitTime * time.Second):
			continue LOOP
		}
	}

	if counter == retryNum {
		err = errors.New("time is over")
	}

	return result, err
}

// Get list of connection addresses.
func (node *Node) Connections() []string {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	var list []string
	for addr := range node.connections {
		list = append(list, addr)
	}

	return list
}

// Check the existence of an address in the list of connections.
func (node *Node) InConnections(address string) bool {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	_, ok := node.connections[address]
	return ok
}

// Connect to node by address.
// Client handle function need be not null.
func (node *Node) Connect(address string) error {
	if node.hasMaxConnSize() {
		return errors.New("max conn")
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	node.setConnection(address, conn)
	go node.handleConn(address)

	return nil
}

// Disconnect from node by address.
func (node *Node) Disconnect(address string) {
	if node.InConnections(address) {
		node.getConnection(address).Close()
	}
	node.delConnection(address)
}

func (node *Node) handleConn(id string) {
	var (
		counter  = uint(0)
		retryNum = gopeer.Get("RETRY_NUM").(uint)
		conn     = node.getConnection(id)
	)

	defer func() {
		conn.Close()
		node.delConnection(id)
	}()

	for {
		if counter == retryNum {
			break
		}

		msg := readMessage(conn)
		if msg == nil {
			counter++
			continue
		}

		counter = 0

	REPEAT:
		if node.inMapping(msg.Body.Hash) {
			continue
		}
		node.setMapping(msg.Body.Hash)

		node.send(msg)

		decMsg := node.client.Decrypt(msg)
		if decMsg == nil {
			continue
		}

		sender := crypto.LoadPubKey(decMsg.Head.Sender)
		if sender == nil {
			continue
		}

		if node.f2f.State() && !node.f2f.InList(sender) {
			continue
		}

		if bytes.Equal(decMsg.Head.Title, gopeer.Get("ROUTE_MSG").([]byte)) {
			msg = local.Package(decMsg.Body.Data).Deserialize()
			goto REPEAT
		}

		node.handleFunc(decMsg)
	}
}

func (node *Node) handleFunc(msg *local.Message) {
	nonceSize := gopeer.Get("SALT_SIZE").(uint)

	if uint(len(msg.Head.Title)) < nonceSize {
		return
	}

	nonce := msg.Head.Title[:nonceSize]
	fname := msg.Head.Title[nonceSize:]

	respBytes := gopeer.Get("RET_BYTES").([]byte)

	// Receive response
	if bytes.HasPrefix(fname, respBytes) {
		node.response(
			nonce,
			msg.Body.Data,
		)
		return
	}

	// Send response
	handler := node.getFunction(fname)
	if handler == nil {
		return
	}

	node.send(node.client.Encrypt(
		crypto.LoadPubKey(msg.Head.Sender),
		local.NewMessage(
			bytes.Join(
				[][]byte{
					nonce,
					respBytes,
					fname,
				},
				[]byte{},
			),
			handler(node.client, msg),
			gopeer.Get("POWS_DIFF").(uint),
		),
	))
}

func (node *Node) send(msg *local.Message) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	pack := msg.Serialize()
	bytesMsg := bytes.Join(
		[][]byte{
			pack.SizeToBytes(),
			pack.Bytes(),
		},
		[]byte{},
	)

	node.mapping[encoding.Base64Encode(msg.Body.Hash)] = true
	for _, cn := range node.connections {
		go cn.Write(bytesMsg)
	}
}

func (node *Node) response(nonce []byte, data []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	strNonce := encoding.Base64Encode(nonce)
	if _, ok := node.actions[strNonce]; ok {
		node.actions[strNonce] <- data
	}
}

func (node *Node) setFunction(name []byte, handle func(*local.Client, *local.Message) []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.hroutes[encoding.Base64Encode(name)] = handle
}

func (node *Node) getFunction(name []byte) func(*local.Client, *local.Message) []byte {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return node.hroutes[encoding.Base64Encode(name)]
}

func (node *Node) setAction(nonce []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.actions[encoding.Base64Encode(nonce)] = make(chan []byte)
}

func (node *Node) getAction(nonce []byte) chan []byte {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	ch, ok := node.actions[encoding.Base64Encode(nonce)]
	if !ok {
		return make(chan []byte)
	}

	return ch
}

func (node *Node) delAction(nonce []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	delete(node.actions, encoding.Base64Encode(nonce))
}

func (node *Node) setMapping(hash []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	if uint(len(node.mapping)) > gopeer.Get("MAPP_SIZE").(uint) {
		node.mapping = make(map[string]bool)
	}

	node.mapping[encoding.Base64Encode(hash)] = true
}

func (node *Node) inMapping(hash []byte) bool {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	_, ok := node.mapping[encoding.Base64Encode(hash)]
	return ok
}

func (node *Node) hasMaxConnSize() bool {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return uint(len(node.connections)) > gopeer.Get("CONN_SIZE").(uint)
}

func (node *Node) setConnection(id string, conn net.Conn) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.connections[id] = conn
}

func (node *Node) getConnection(id string) net.Conn {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return node.connections[id]
}

func (node *Node) delConnection(id string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	delete(node.connections, id)
}

type Nonce []byte

func (nonce Nonce) Cmp(slice []byte) int {
	return bytes.Compare(nonce, slice)
}
