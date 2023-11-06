package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex        sync.Mutex
	fListener     net.Listener
	fSettings     ISettings
	fHashMapping  map[string]struct{}
	fConnections  map[string]conn.IConn
	fHandleRoutes map[uint64]IHandlerF
}

// Creating a node object managed by connections with multiple nodes.
// Saves hashes of received messages to a buffer to prevent network cycling.
// Redirects messages to handle routers by keys.
func NewNode(pSett ISettings) INode {
	return &sNode{
		fSettings:     pSett,
		fHashMapping:  make(map[string]struct{}, pSett.GetCapacity()),
		fConnections:  make(map[string]conn.IConn),
		fHandleRoutes: make(map[uint64]IHandlerF),
	}
}

// Return settings interface.
func (p *sNode) GetSettings() ISettings {
	return p.fSettings
}

// Puts the hash of the message in the buffer and sends the message to all connections of the node.
func (p *sNode) BroadcastMessage(pMsg message.IMessage) error {
	_ = p.inMappingWithSet(pMsg.GetHash()) // node can redirect received message

	listErr := make([]error, 0, p.fSettings.GetMaxConnects())

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for a, c := range p.GetConnections() {
		wg.Add(1)

		chErr := make(chan error)
		go func(c conn.IConn) {
			chErr <- c.WriteMessage(pMsg)
		}(c)

		go func(a string, c conn.IConn) {
			var resErr error

			defer wg.Done()
			defer func() {
				if resErr != nil {
					_ = p.DelConnection(a)
				}

				mutex.Lock()
				listErr = append(listErr, resErr)
				mutex.Unlock()
			}()

			select {
			case err := <-chErr:
				resErr = err // err can be = nil
			case <-time.After(p.fSettings.GetWriteTimeout()):
				<-chErr
				resErr = errors.NewError(
					fmt.Sprintf(
						"write timeout %s",
						c.GetSocket().RemoteAddr().String(),
					),
				)
			}
		}(a, c)
	}

	wg.Wait()

	var resErr error
	for _, err := range listErr {
		resErr = errors.AppendError(resErr, err)
	}
	return resErr
}

// Opens a tcp connection to receive data from outside.
// Checks the number of valid connections.
// Redirects connections to the handle router.
func (p *sNode) Run() error {
	listener, err := net.Listen("tcp", p.fSettings.GetAddress())
	if err != nil {
		return errors.WrapError(err, "run node")
	}

	go func(pListener net.Listener) {
		defer pListener.Close()
		p.setListener(pListener)
		for {
			tconn, err := p.getListener().Accept()
			if err != nil {
				break
			}

			if p.hasMaxConnSize() {
				tconn.Close()
				continue
			}

			sett := p.fSettings.GetConnSettings()
			conn := conn.LoadConn(sett, tconn)
			address := tconn.RemoteAddr().String()

			p.setConnection(address, conn)
			go p.handleConn(address, conn)
		}
	}(listener)

	return nil
}

// Closes the listener and all connections.
func (p *sNode) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	toClose := make([]types.ICloser, 0, len(p.fConnections)+1)
	if p.fListener != nil {
		toClose = append(toClose, p.fListener)
	}

	for id, conn := range p.fConnections {
		toClose = append(toClose, conn)
		delete(p.fConnections, id)
	}

	if err := types.CloseAll(toClose); err != nil {
		return errors.WrapError(err, "stop node")
	}
	return nil
}

// Saves the function to the map by key for subsequent redirection.
func (p *sNode) HandleFunc(pHead uint64, pHandle IHandlerF) INode {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fHandleRoutes[pHead] = pHandle
	return p
}

// Retrieves the entire list of connections with addresses.
func (p *sNode) GetConnections() map[string]conn.IConn {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	var mapping = make(map[string]conn.IConn, len(p.fConnections))
	for addr, conn := range p.fConnections {
		mapping[addr] = conn
	}

	return mapping
}

// Connects to the node at the specified address and automatically starts reading all incoming messages.
// Checks the number of connections.
func (p *sNode) AddConnection(pAddress string) error {
	if p.hasMaxConnSize() {
		return errors.NewError("has max connections size")
	}

	if _, ok := p.getConnection(pAddress); ok {
		return errors.NewError("connection already exist")
	}

	sett := p.fSettings.GetConnSettings()
	conn, err := conn.NewConn(sett, pAddress)
	if err != nil {
		return errors.WrapError(err, "add connect")
	}

	p.setConnection(pAddress, conn)
	go p.handleConn(pAddress, conn)

	return nil
}

// Disables the connection at the address and removes the connection from the connection list.
func (p *sNode) DelConnection(pAddress string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	conn, ok := p.fConnections[pAddress]
	if !ok {
		return errors.NewError("unknown connect")
	}

	delete(p.fConnections, pAddress)

	if err := conn.Close(); err != nil {
		return errors.WrapError(err, "connect close")
	}

	return nil
}

// Processes the received data from the connection.
func (p *sNode) handleConn(pAddress string, pConn conn.IConn) {
	defer p.DelConnection(pAddress)
	for {
		var (
			readerCh = make(chan struct{})
			returnCh = make(chan bool)
		)

		go func() {
			msg, err := pConn.ReadMessage(readerCh)
			if err != nil {
				returnCh <- false
				return
			}
			returnCh <- p.handleMessage(pConn, msg)
		}()

		<-readerCh
		select {
		case ok := <-returnCh:
			if !ok {
				return
			}
		case <-time.After(p.fSettings.GetReadTimeout()):
			<-returnCh
			return
		}
	}
}

// Processes the message for correctness and redirects it to the handler function.
// Returns true if the message was successfully redirected to the handler function
// > or if the message already existed in the hash value store.
func (p *sNode) handleMessage(pConn conn.IConn, pMsg message.IMessage) bool {
	// check message in mapping by hash
	if p.inMappingWithSet(pMsg.GetHash()) {
		return true
	}

	// get function by head
	pld := pMsg.GetPayload()
	f, ok := p.getFunction(pld.GetHead())
	if !ok || f == nil {
		// function is not found = protocol error
		return false
	}

	if err := f(p, pConn, pMsg); err != nil {
		// function error = protocol error
		return false
	}

	return true
}

// Checks the current number of connections with the limit.
func (p *sNode) hasMaxConnSize() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	maxConns := p.fSettings.GetMaxConnects()
	return uint64(len(p.fConnections)) >= maxConns
}

// Checks the hash of the message for existence in the hash store.
// Returns true if the hash already existed, otherwise false.
func (p *sNode) inMappingWithSet(pHash []byte) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	sHash := encoding.HexEncode(pHash)

	// skey already exists
	if _, ok := p.fHashMapping[sHash]; ok {
		return true
	}

	// push skey to mapping
	p.fHashMapping[sHash] = struct{}{}
	return false
}

// Saves the connection to the map.
func (p *sNode) getConnection(pAddress string) (conn.IConn, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	conn, ok := p.fConnections[pAddress]
	return conn, ok
}

// Saves the connection to the map.
func (p *sNode) setConnection(pAddress string, pConn conn.IConn) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fConnections[pAddress] = pConn
}

// Gets the handler function by key.
func (p *sNode) getFunction(pHead uint64) (IHandlerF, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	f, ok := p.fHandleRoutes[pHead]
	return f, ok
}

// Sets the listener.
func (p *sNode) setListener(pListener net.Listener) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fListener = pListener
}

// Gets the listener.
func (p *sNode) getListener() net.Listener {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fListener
}
