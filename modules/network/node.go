package network

import (
	"net"
	"sync"

	"github.com/number571/go-peer/modules/payload"
	"github.com/number571/go-peer/modules/storage"
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
	fConnections  map[string]IConn
	fHandleRoutes map[uint64]IHandlerF
}

// Create client by private key as identification.
func NewNode(sett ISettings) INode {
	return &sNode{
		fSettings:     sett,
		fHashMapping:  storage.NewMemoryStorage(sett.GetCapacity()),
		fConnections:  make(map[string]IConn),
		fHandleRoutes: make(map[uint64]IHandlerF),
	}
}

func (node *sNode) Broadcast(pl payload.IPayload) error {
	// set this message to mapping
	msg := NewMessage(pl)
	node.inMappingWithSet(msg.Hash())

	var err error
	for _, conn := range node.Connections() {
		e := conn.Write(msg)
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
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listen.Close()

	node.fListener = listen
	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}

		node.fMutex.Lock()
		isConnLimit := node.hasMaxConnSize()
		node.fMutex.Unlock()
		if isConnLimit {
			conn.Close()
			continue
		}

		node.fMutex.Lock()
		iconn := LoadConn(node.fSettings, conn)
		address := iconn.Socket().RemoteAddr().String()
		node.fConnections[address] = iconn
		node.fMutex.Unlock()

		go node.handleConn(address, iconn)
	}

	return nil
}

func (node *sNode) Close() error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	var err error

	for id, conn := range node.fConnections {
		e := conn.Close()
		if e != nil {
			err = e
		}
		delete(node.fConnections, id)
	}
	if node.fListener != nil {
		e := node.fListener.Close()
		if e != nil {
			err = e
		}
	}

	return err
}

// Add function to mapping for route use.
func (node *sNode) Handle(head uint64, handle IHandlerF) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fHandleRoutes[head] = handle
	return node
}

func (node *sNode) handleConn(address string, conn IConn) {
	defer node.Disconnect(address)
	for {
		ok := node.handleMessage(conn, conn.Read())
		if !ok {
			break
		}
	}
}

// Get list of connection addresses.
func (node *sNode) Connections() map[string]IConn {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	var mapping = make(map[string]IConn)
	for addr, conn := range node.fConnections {
		mapping[addr] = conn
	}

	return mapping
}

// Connect to node by address.
// Client handle function need be not null.
func (node *sNode) Connect(address string) IConn {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	if node.hasMaxConnSize() {
		return nil
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}

	iconn := LoadConn(node.fSettings, conn)
	node.fConnections[address] = iconn

	go node.handleConn(address, iconn)
	return iconn
}

func (node *sNode) Disconnect(address string) error {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	var err error
	conn, ok := node.fConnections[address]
	if ok {
		err = conn.Close()
	}

	delete(node.fConnections, address)
	return err
}

func (node *sNode) handleMessage(conn IConn, msg IMessage) bool {
	// null message from connection is error
	if msg == nil {
		return false
	}

	// check message in mapping by hash
	if node.inMappingWithSet(msg.Hash()) {
		return true
	}

	// get function by head
	f, ok := node.getFunction(msg.Payload().Head())
	if !ok || f == nil {
		return false
	}

	f(node, conn, msg.Payload())
	return true
}

func (node *sNode) hasMaxConnSize() bool {
	maxConns := node.fSettings.GetMaxConnects()
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

func (node *sNode) getFunction(head uint64) (IHandlerF, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	f, ok := node.fHandleRoutes[head]
	return f, ok
}
