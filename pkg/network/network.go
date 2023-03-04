package network

import (
	"fmt"
	"net"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex        sync.Mutex
	fListener     net.Listener
	fSettings     ISettings
	fHashMapping  storage.IKeyValueStorage
	fConnections  map[string]conn.IConn
	fHandleRoutes map[uint64]IHandlerF
}

// Creating a node object managed by connections with multiple nodes.
// Saves hashes of received messages to a buffer to prevent network cycling.
// Redirects messages to handle routers by keys.
func NewNode(sett ISettings) INode {
	return &sNode{
		fSettings:     sett,
		fHashMapping:  storage.NewMemoryStorage(sett.GetCapacity()),
		fConnections:  make(map[string]conn.IConn),
		fHandleRoutes: make(map[uint64]IHandlerF),
	}
}

// Return settings interface.
func (node *sNode) GetSettings() ISettings {
	return node.fSettings
}

// Puts the hash of the message in the buffer and sends the message to all connections of the node.
func (node *sNode) BroadcastPayload(pld payload.IPayload) error {
	hasher := hashing.NewSHA256Hasher(pld.ToBytes())
	node.inMappingWithSet(hasher.ToBytes())

	var err error
	for _, conn := range node.GetConnections() {
		e := conn.WritePayload(pld)
		if e != nil {
			err = e
		}
	}

	return err
}

// Opens a tcp connection to receive data from outside.
// Checks the number of valid connections.
// Redirects connections to the handle router.
func (node *sNode) Run() error {
	listener, err := net.Listen("tcp", node.GetSettings().GetAddress())
	if err != nil {
		return err
	}
	defer listener.Close()

	node.setListener(listener)
	for {
		tconn, err := node.getListener().Accept()
		if err != nil {
			break
		}

		if node.hasMaxConnSize() {
			tconn.Close()
			continue
		}

		sett := node.GetSettings().GetConnSettings()
		conn := conn.LoadConn(sett, tconn)
		address := tconn.RemoteAddr().String()

		node.setConnection(address, conn)
		go node.handleConn(address, conn)
	}

	return nil
}

// Closes the listener and all connections.
func (node *sNode) Stop() error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	toClose := make([]types.ICloser, 0, len(node.fConnections)+1)
	if node.fListener != nil {
		toClose = append(toClose, node.fListener)
	}

	for id, conn := range node.fConnections {
		toClose = append(toClose, conn)
		delete(node.fConnections, id)
	}

	return types.CloseAll(toClose)
}

// Saves the function to the map by key for subsequent redirection.
func (node *sNode) HandleFunc(head uint64, handle IHandlerF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
	return node
}

// Retrieves the entire list of connections with addresses.
func (node *sNode) GetConnections() map[string]conn.IConn {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	var mapping = make(map[string]conn.IConn, len(node.fConnections))
	for addr, conn := range node.fConnections {
		mapping[addr] = conn
	}

	return mapping
}

// Connects to the node at the specified address and automatically starts reading all incoming messages.
// Checks the number of connections.
func (node *sNode) AddConnect(address string) error {
	if node.hasMaxConnSize() {
		return fmt.Errorf("has max connections size")
	}

	sett := node.GetSettings().GetConnSettings()
	conn, err := conn.NewConn(sett, address)
	if err != nil {
		return err
	}

	node.setConnection(address, conn)
	go node.handleConn(address, conn)

	return nil
}

// Disables the connection at the address and removes the connection from the connection list.
func (node *sNode) DelConnect(address string) error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	conn, ok := node.fConnections[address]
	if !ok {
		return nil
	}

	delete(node.fConnections, address)
	return conn.Close()
}

// Processes the received data from the connection.
func (node *sNode) handleConn(address string, conn conn.IConn) {
	defer node.DelConnect(address)
	for {
		ok := node.handleMessage(conn, conn.ReadPayload())
		if !ok {
			break
		}
	}
}

// Processes the message for correctness and redirects it to the handler function.
// Returns true if the message was successfully redirected to the handler function
// > or if the message already existed in the hash value store.
func (node *sNode) handleMessage(conn conn.IConn, pld payload.IPayload) bool {
	// null message from connection is error
	if pld == nil {
		return false
	}

	// check message in mapping by hash
	hash := hashing.NewSHA256Hasher(pld.ToBytes()).ToBytes()
	if node.inMappingWithSet(hash) {
		return true
	}

	// get function by head
	f, ok := node.getFunction(pld.GetHead())
	if !ok || f == nil {
		return false
	}

	f(node, conn, pld.GetBody())
	return true
}

// Checks the current number of connections with the limit.
func (node *sNode) hasMaxConnSize() bool {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	maxConns := node.GetSettings().GetMaxConnects()
	return uint64(len(node.fConnections)) > maxConns
}

// Checks the hash of the message for existence in the hash store.
// Returns true if the hash already existed, otherwise false.
func (node *sNode) inMappingWithSet(hash []byte) bool {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	// skey already exists
	_, err := node.fHashMapping.Get(hash)
	if err == nil {
		return true
	}

	// push skey to mapping
	node.fHashMapping.Set(hash, []byte{1})
	return false
}

// Saves the connection to the map.
func (node *sNode) setConnection(address string, conn conn.IConn) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fConnections[address] = conn
}

// Gets the handler function by key.
func (node *sNode) getFunction(head uint64) (IHandlerF, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleRoutes[head]
	return f, ok
}

// Sets the listener.
func (node *sNode) setListener(listener net.Listener) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fListener = listener
}

// Gets the listener.
func (node *sNode) getListener() net.Listener {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fListener
}
