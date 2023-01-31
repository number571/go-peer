package network

import (
	"fmt"
	"net"
	"sync"

	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ INode = &sNode{}
)

// Basic structure for network use.
type sNode struct {
	fMutex        sync.Mutex
	fListener     net.Listener
	fSettings     ISettings
	fHashMapping  storage.IKeyValueStorage
	fConnections  map[string]conn.IConn
	fHandleRoutes map[uint64]IHandlerF
}

// Create client by private key as identification.
func NewNode(sett ISettings) INode {
	return &sNode{
		fSettings:     sett,
		fHashMapping:  storage.NewMemoryStorage(sett.GetCapacity()),
		fConnections:  make(map[string]conn.IConn),
		fHandleRoutes: make(map[uint64]IHandlerF),
	}
}

func (node *sNode) Broadcast(pld payload.IPayload) error {
	// set this message to mapping
	hash := hashing.NewSHA256Hasher(pld.Bytes()).Bytes()
	node.inMappingWithSet(hash)

	var err error
	for _, conn := range node.Connections() {
		e := conn.Write(pld)
		if e != nil {
			err = e
		}
	}

	return err
}

func (node *sNode) Settings() ISettings {
	return node.fSettings
}

// Turn on listener by address.
// Client handle function need be not null.
func (node *sNode) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
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

		sett := node.Settings().GetConnSettings()
		conn := conn.LoadConn(sett, tconn)
		address := tconn.RemoteAddr().String()

		node.setConnection(address, conn)
		go node.handleConn(address, conn)
	}

	return nil
}

func (node *sNode) Close() error {
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

	return closer.CloseAll(toClose)
}

// Add function to mapping for route use.
func (node *sNode) Handle(head uint64, handle IHandlerF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
	return node
}

func (node *sNode) handleConn(address string, conn conn.IConn) {
	defer node.Disconnect(address)
	for {
		ok := node.handleMessage(conn, conn.Read())
		if !ok {
			break
		}
	}
}

// Get list of connection addresses.
func (node *sNode) Connections() map[string]conn.IConn {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	var mapping = make(map[string]conn.IConn, len(node.fConnections))
	for addr, conn := range node.fConnections {
		mapping[addr] = conn
	}

	return mapping
}

// Connect to node by address.
// Client handle function need be not null.
func (node *sNode) Connect(address string) (conn.IConn, error) {
	if node.hasMaxConnSize() {
		return nil, fmt.Errorf("has max connections size")
	}

	sett := node.Settings().GetConnSettings()
	conn, err := conn.NewConn(sett, address)
	if err != nil {
		return nil, err
	}

	node.setConnection(address, conn)
	go node.handleConn(address, conn)

	return conn, nil
}

func (node *sNode) Disconnect(address string) error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	conn, ok := node.fConnections[address]
	if !ok {
		return nil
	}

	delete(node.fConnections, address)
	return conn.Close()
}

func (node *sNode) handleMessage(conn conn.IConn, pld payload.IPayload) bool {
	// null message from connection is error
	if pld == nil {
		return false
	}

	// check message in mapping by hash
	hash := hashing.NewSHA256Hasher(pld.Bytes()).Bytes()
	if node.inMappingWithSet(hash) {
		return true
	}

	// get function by head
	f, ok := node.getFunction(pld.Head())
	if !ok || f == nil {
		return false
	}

	f(node, conn, pld.Body())
	return true
}

func (node *sNode) hasMaxConnSize() bool {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	maxConns := node.Settings().GetMaxConnects()
	return uint64(len(node.fConnections)) > maxConns
}

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

func (node *sNode) setConnection(address string, conn conn.IConn) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fConnections[address] = conn
}

func (node *sNode) getFunction(head uint64) (IHandlerF, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleRoutes[head]
	return f, ok
}

func (node *sNode) setListener(listener net.Listener) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fListener = listener
}

func (node *sNode) getListener() net.Listener {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fListener
}
