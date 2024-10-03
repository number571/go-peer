package network

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/utils"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ INode = &sNode{}
)

type sNode struct {
	fMutex        sync.RWMutex
	fSettings     ISettings
	fListener     net.Listener
	fCacheSetter  cache.ICacheSetter
	fConnections  map[string]conn.IConn
	fHandleRoutes map[uint32]IHandlerF
}

// Creating a node object managed by connections with multiple nodes.
// Saves hashes of received messages to a buffer to prevent network cycling.
// Redirects messages to handle routers by keys.
func NewNode(
	pSettings ISettings,
	pCacheSetter cache.ICacheSetter,
) INode {
	return &sNode{
		fSettings:     pSettings,
		fCacheSetter:  pCacheSetter,
		fConnections:  make(map[string]conn.IConn, pSettings.GetMaxConnects()),
		fHandleRoutes: make(map[uint32]IHandlerF, 64),
	}
}

// Return settings interface.
func (p *sNode) GetSettings() ISettings {
	return p.fSettings
}

// Puts the hash of the message in the buffer and sends the message to all connections of the node.
func (p *sNode) BroadcastMessage(pCtx context.Context, pMsg net_message.IMessage) error {
	connections := p.GetConnections()
	lenConnections := len(connections)

	// can't broadcast message to the network if len(connections) = 0
	if lenConnections == 0 {
		return ErrNoConnections
	}

	// node can redirect received message
	_ = p.fCacheSetter.Set(pMsg.GetHash(), []byte{})

	wg := sync.WaitGroup{}
	wg.Add(lenConnections)

	listErr := make([]error, lenConnections)
	i := 0

	for a, c := range connections {
		chErr := make(chan error)

		go func(c conn.IConn) {
			chErr <- c.WriteMessage(pCtx, pMsg)
		}(c)

		go func(i int, a string) {
			defer wg.Done()

			timer := time.NewTimer(p.fSettings.GetWriteTimeout())
			defer timer.Stop()

			select {
			case <-pCtx.Done():
				listErr[i] = pCtx.Err()
			case <-timer.C:
				listErr[i] = ErrWriteTimeout
			case err := <-chErr:
				if err == nil {
					return
				}
				listErr[i] = utils.MergeErrors(ErrBroadcastMessage, err)
			}

			// if got error -> delete connection
			_ = p.DelConnection(a)
		}(i, a)

		i++
	}

	wg.Wait()
	return utils.MergeErrors(listErr...)
}

// Opens a tcp connection to receive data from outside.
// Checks the number of valid connections.
// Redirects connections to the handle router.
func (p *sNode) Listen(pCtx context.Context) error {
	listener, err := net.Listen("tcp", p.fSettings.GetAddress())
	if err != nil {
		return utils.MergeErrors(ErrCreateListener, err)
	}
	defer listener.Close()

	p.setListener(listener)
	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			tconn, err := p.getListener().Accept()
			if err != nil {
				return utils.MergeErrors(ErrListenerAccept, err)
			}

			if p.hasMaxConnSize() {
				tconn.Close()
				continue
			}

			conn := conn.LoadConn(p.fSettings.GetConnSettings(), tconn)
			address := tconn.RemoteAddr().String()

			p.setConnection(address, conn)
			go p.handleConn(pCtx, address, conn)
		}
	}
}

// Closes the listener and all connections.
func (p *sNode) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	listErr := make([]error, 0, len(p.fConnections)+1)
	if p.fListener != nil {
		listErr = append(listErr, p.fListener.Close())
	}

	for id, conn := range p.fConnections {
		delete(p.fConnections, id)
		listErr = append(listErr, conn.Close())
	}

	return utils.MergeErrors(listErr...)
}

// Saves the function to the map by key for subsequent redirection.
func (p *sNode) HandleFunc(pHead uint32, pHandle IHandlerF) INode {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fHandleRoutes[pHead] = pHandle
	return p
}

// Retrieves the entire list of connections with addresses.
func (p *sNode) GetConnections() map[string]conn.IConn {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	mapping := make(map[string]conn.IConn, len(p.fConnections))
	for addr, conn := range p.fConnections {
		mapping[addr] = conn
	}

	return mapping
}

// Connects to the node at the specified address and automatically starts reading all incoming messages.
// Checks the number of connections.
func (p *sNode) AddConnection(pCtx context.Context, pAddress string) error {
	if p.hasMaxConnSize() {
		return ErrHasLimitConnections
	}

	if _, ok := p.getConnection(pAddress); ok {
		return ErrConnectionIsExist
	}

	sett := p.fSettings.GetConnSettings()
	conn, err := conn.Connect(pCtx, sett, pAddress)
	if err != nil {
		return utils.MergeErrors(ErrAddConnections, err)
	}

	p.setConnection(pAddress, conn)
	go p.handleConn(pCtx, pAddress, conn)

	return nil
}

// Disables the connection at the address and removes the connection from the connection list.
func (p *sNode) DelConnection(pAddress string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	conn, ok := p.fConnections[pAddress]
	if !ok {
		return ErrConnectionIsNotExist
	}

	delete(p.fConnections, pAddress)

	if err := conn.Close(); err != nil {
		return utils.MergeErrors(ErrCloseConnection, err)
	}

	return nil
}

// Processes the received data from the connection.
func (p *sNode) handleConn(pCtx context.Context, pAddress string, pConn conn.IConn) {
	defer func() { _ = p.DelConnection(pAddress) }()

	var (
		readHeadCh = make(chan struct{})
		readFullCh = make(chan net_message.IMessage)
	)

	go p.messageReader(
		pCtx,
		pConn,
		readHeadCh,
		readFullCh,
	)

	for {
		select {
		case <-pCtx.Done():
			return
		case <-readHeadCh:
			select {
			case <-pCtx.Done():
				return
			case <-time.After(p.fSettings.GetReadTimeout()):
				return
			case msg := <-readFullCh:
				if msg == nil {
					return
				}
				if ok := p.handleMessage(pCtx, pConn, msg); !ok {
					return
				}
				break
			}
		}
	}
}

func (p *sNode) messageReader(
	pCtx context.Context,
	pConn conn.IConn,
	readHeadCh chan<- struct{},
	readFullCh chan<- net_message.IMessage,
) {
	for {
		select {
		case <-pCtx.Done():
			return
		default:
			msg, err := pConn.ReadMessage(pCtx, readHeadCh)
			if err != nil {
				readFullCh <- nil
				return
			}
			readFullCh <- msg
		}
	}
}

// Processes the message for correctness and redirects it to the handler function.
// Returns true if the message was successfully redirected to the handler function
// > or if the message already existed in the hash value store.
func (p *sNode) handleMessage(pCtx context.Context, pConn conn.IConn, pMsg net_message.IMessage) bool {
	if !p.fCacheSetter.Set(pMsg.GetHash(), []byte{}) {
		return true // hash of message already in queue
	}

	f, ok := p.getFunction(pMsg.GetPayload().GetHead())
	if !ok || f == nil {
		return false // function is not found = protocol error
	}

	err := f(pCtx, p, pConn, pMsg)
	return err == nil // function error = protocol error
}

// Checks the current number of connections with the limit.
func (p *sNode) hasMaxConnSize() bool {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	maxConns := p.fSettings.GetMaxConnects()
	return uint64(len(p.fConnections)) >= maxConns
}

// Saves the connection to the map.
func (p *sNode) getConnection(pAddress string) (conn.IConn, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

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
func (p *sNode) getFunction(pHead uint32) (IHandlerF, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

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
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fListener
}
