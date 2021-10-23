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
// Handle function is used when the network exists. Can be null.
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

		if node.isMaxConnSize() {
			conn.Close()
			continue
		}

		id := encoding.Base64Encode(crypto.Rand(gopeer.Get("RAND_SIZE").(uint)))
		node.setConnection(id, conn)

		go node.handleConn(id)
	}

	return nil
}

// Add function to mapping for route use.
func (node *Node) Handle(title []byte, handle func(*local.Client, *local.Message) []byte) *Node {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.hroutes[encoding.Base64Encode(title)] = handle
	return node
}

// Send message by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (node *Node) Send(msg *local.Message, route *local.Route) ([]byte, error) {
	var (
		err      error
		result   []byte
		hash     = route.Receiver().Address()
		retryNum = gopeer.Get("RETRY_NUM").(uint)
	)

	node.setAction(hash)
	defer node.delAction(hash)

REPEAT:
	routeMsg := node.client.RouteMessage(msg, route)
	if routeMsg == nil {
		return result, errors.New("psender is nil and routes not nil")
	}

	node.send(routeMsg)
	waitTime := time.Duration(gopeer.Get("WAIT_TIME").(uint))

	select {
	case result = <-node.getAction(hash):
	case <-time.After(waitTime * time.Second):
		if retryNum > 1 {
			retryNum -= 1
			goto REPEAT
		}
		err = errors.New("time is over")
	}

	node.delAction(hash)
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
	if node.isMaxConnSize() {
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
	const (
		NUM_FAILS = 5
	)

	var (
		failCounter = 0
		conn        = node.getConnection(id)
	)

	defer func() {
		conn.Close()
		node.delConnection(id)
	}()

	for {
		if failCounter == NUM_FAILS {
			break
		}

		msg := readMessage(conn)
		if msg == nil {
			failCounter++
			continue
		}

		failCounter = 0

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
	fname := msg.Head.Title
	respBytes := gopeer.Get("RET_BYTES").([]byte)

	// Receive response
	if bytes.HasPrefix(fname, respBytes) {
		node.response(
			crypto.LoadPubKey(msg.Head.Sender),
			msg.Body.Data,
		)
		return
	}

	// Send response
	diff := gopeer.Get("POWS_DIFF").(uint)

	handler := node.getFunction(fname)
	if handler == nil {
		return
	}

	node.send(node.client.Encrypt(
		crypto.LoadPubKey(msg.Head.Sender),
		local.NewMessage(
			bytes.Join([][]byte{
				gopeer.Get("RET_BYTES").([]byte),
				fname,
			}, []byte{}),
			handler(node.client, msg),
			diff,
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

func (node *Node) response(pub crypto.PubKey, data []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	hash := pub.Address()
	if _, ok := node.actions[hash]; ok {
		node.actions[hash] <- data
	}
}

func (node *Node) getFunction(name []byte) func(*local.Client, *local.Message) []byte {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return node.hroutes[encoding.Base64Encode(name)]
}

func (node *Node) setAction(hash string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.actions[hash] = make(chan []byte)
}

func (node *Node) getAction(hash string) chan []byte {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	ch, ok := node.actions[hash]
	if !ok {
		return make(chan []byte)
	}

	return ch
}

func (node *Node) delAction(hash string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	delete(node.actions, hash)
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

func (node *Node) isMaxConnSize() bool {
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
