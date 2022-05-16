package network

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
)

var (
	_ INode = &sNode{}
)

// Basic structure for network use.
type sNode struct {
	fMutex       sync.Mutex
	fListener    net.Listener
	fClient      local.IClient
	fHRoutes     map[string]iHandler
	fMapping     map[string]bool
	fConnections map[string]net.Conn
	fF2F         iF2F
	fChecker     iChecker
	fPseudo      iPseudo
	fOnline      iOnline
	fRouter      iRouter
	fActions     map[string]chan []byte
}

// Create client by private key as identification.
func NewNode(client local.IClient) INode {
	if client == nil {
		return nil
	}

	node := &sNode{
		fClient:      client,
		fHRoutes:     make(map[string]iHandler),
		fMapping:     make(map[string]bool),
		fConnections: make(map[string]net.Conn),
		fF2F: &sF2F{
			fMapping: make(map[string]crypto.IPubKey),
		},
		fChecker: &sChecker{
			fChannel: make(chan struct{}),
			fMapping: make(map[string]iCheckerInfo),
		},
		fPseudo: &sPseudo{
			fChannel: make(chan struct{}),
			fPrivKey: crypto.NewPrivKey(client.PubKey().Size()),
		},
		fOnline:  &sOnline{},
		fRouter:  func(_ INode) []crypto.IPubKey { return nil },
		fActions: make(map[string]chan []byte),
	}

	// recurrent structures
	{
		checker := node.fChecker.(*sChecker)
		pseudo := node.fPseudo.(*sPseudo)
		online := node.fOnline.(*sOnline)

		checker.fNode = node
		pseudo.fNode = node
		online.fNode = node
	}

	sett := node.Client().Settings()
	patt := encoding.Uint64ToBytes(sett.Get(settings.MaskPing))

	return node.Handle(patt, nil)
}

func (node *sNode) WithResponseRouter(router iRouter) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fRouter = router
	return node
}

// Close checker, pseudo, online status, listener and current connections.
func (node *sNode) Close() {
	statuses := []iStatus{
		node.fChecker,
		node.fPseudo,
		node.fOnline,
	}
	for _, status := range statuses {
		status.Switch(false)
	}

	node.fMutex.Lock()
	if node.fListener != nil {
		node.fListener.Close()
	}
	for id, conn := range node.fConnections {
		conn.Close()
		delete(node.fConnections, id)
	}
	node.fMutex.Unlock()
}

// Return client structure.
func (node *sNode) Client() local.IClient {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fClient
}

// Return checker structure.
func (node *sNode) Checker() iChecker {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fChecker
}

// Return pseudo structure.
func (node *sNode) Pseudo() iPseudo {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fPseudo
}

// Return online structure.
func (node *sNode) Online() iOnline {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fOnline
}

// Return f2f structure.
func (node *sNode) F2F() iF2F {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return node.fF2F
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

		if node.hasMaxConnSize() {
			conn.Close()
			continue
		}

		rsize := node.Client().Settings().Get(settings.SizeSkey)
		id := crypto.NewPRNG().String(rsize)

		node.setConnection(id, conn)
		go node.handleConn(id)
	}

	return nil
}

// Add function to mapping for route use.
func (node *sNode) Handle(title []byte, handle iHandler) INode {
	return node.setFunction(title, handle)
}

// Send message by public key of receiver.
func (node *sNode) Request(route local.IRoute, msg local.IMessage) ([]byte, error) {
	return node.doRequest(
		route,
		msg,
		node.Client().Settings().Get(settings.SizeRtry),
		node.Client().Settings().Get(settings.TimeWait),
	)
}

// Get list of connection addresses.
func (node *sNode) Connections() []string {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	var list []string
	for addr := range node.fConnections {
		list = append(list, addr)
	}

	return list
}

// Check the existence of an address in the list of connections.
func (node *sNode) InConnections(address string) bool {
	_, ok := node.getConnection(address)
	return ok
}

// Connect to node by address.
// Client handle function need be not null.
func (node *sNode) Connect(address string) error {
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
func (node *sNode) Disconnect(address string) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	conn, ok := node.fConnections[address]
	if ok {
		conn.Close()
	}

	delete(node.fConnections, address)
}

func (node *sNode) handleConn(id string) {
	defer node.Disconnect(id)

	var (
		retryNum = node.Client().Settings().Get(settings.SizeRtry)
		conn, _  = node.getConnection(id)
	)

	counter := uint64(0)
	for {
		if counter > retryNum {
			break
		}

		ok := node.handleMessage(node.readMessage(conn))
		if !ok {
			counter++
			continue
		}

		counter = 0
	}
}

func (node *sNode) handleMessage(msg local.IMessage) bool {
	// null message from connection is error
	if msg == nil {
		return false
	}

	// check message in mapping by hash
	if node.inMapping(msg.Body().Hash()) {
		return true
	}
	node.setMapping(msg.Body().Hash())

	// redirect this message to connections
	node.send(msg)

	// try decrypt message
	decMsg, title := node.Client().Decrypt(msg)
	if decMsg == nil {
		return true
	}

	// if this message is just route message
	// then try procedures again
	routeMsg := node.Client().Settings().Get(settings.MaskRout)

	// if is route package then
	// 1/2 generate new pseudo-package and sleep rand time
	// unpack and send new version of package
	if bytes.Equal(title, encoding.Uint64ToBytes(routeMsg)) {
		if node.fPseudo.Status() && crypto.NewPRNG().Bool() {
			// send pseudo message with random sleep
			size := len(decMsg.Body().Data())
			node.fPseudo.Request(size).Sleep()
		}
		// recursive unpack message
		msg = local.LoadPackage(decMsg.Body().Data()).ToMessage()
		return node.handleMessage(msg)
	}

	// sleep random milliseconds
	if node.fPseudo.Status() {
		node.fPseudo.Sleep()
	}

	// if mode is friend-to-friend and sender not in list of f2f
	// then pass this request
	sender := crypto.LoadPubKey(decMsg.Head().Sender())
	if node.fF2F.Status() && !node.fF2F.InList(sender) {
		return true
	}

	// send message to handler
	node.handleFunc(decMsg, title)
	return true
}

func (node *sNode) handleFunc(msg local.IMessage, title []byte) {
	var (
		skeySize  = node.Client().Settings().Get(settings.SizeSkey)
		respNum   = node.Client().Settings().Get(settings.MaskRout)
		respBytes = encoding.Uint64ToBytes(respNum)
	)

	// receive response
	if bytes.HasPrefix(title, respBytes) {
		title = title[len(respBytes):]
		if uint64(len(title)) < skeySize {
			return
		}
		node.response(
			title[:skeySize],
			msg.Body().Data(),
		)
		return
	}

	// send response
	f := node.getFunction(title)
	if f == nil {
		return
	}

	rmsg, _ := node.Client().Encrypt(
		local.NewRoute(crypto.LoadPubKey(msg.Head().Sender())).
			WithRedirects(node.fPseudo.PrivKey(), node.fRouter(node)),
		local.NewMessage(
			bytes.Join(
				[][]byte{
					respBytes,
					msg.Head().Session(),
					title,
				},
				[]byte{},
			),
			f(node, msg),
		),
	)
	node.send(rmsg)
}

// Request with retry number and time out.
func (node *sNode) doRequest(route local.IRoute, msg local.IMessage, retryNum, timeOut uint64) ([]byte, error) {
	if len(node.Connections()) == 0 {
		return nil, errors.New("length of connections = 0")
	}

	routeMsg, session := node.Client().Encrypt(route, msg)
	if routeMsg == nil {
		return nil, errors.New("psender is nil and routes not nil")
	}

	node.setAction(session)
	defer node.delAction(session)

	for counter := uint64(0); counter <= retryNum; counter++ {
		node.send(routeMsg)
		resp, err := node.recv(session, timeOut)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			continue
		}
		return resp, nil
	}

	return nil, errors.New("time is over")
}

func (node *sNode) recv(session []byte, timeOut uint64) ([]byte, error) {
	select {
	case result, opened := <-node.getAction(session):
		if !opened {
			return nil, errors.New("chan is closed")
		}
		return result, nil
	case <-time.After(time.Duration(timeOut) * time.Second):
		return nil, nil
	}
}

func (node *sNode) send(msg local.IMessage) {
	node.fMutex.Lock() // TODO
	defer node.fMutex.Unlock()

	pack := msg.ToPackage()
	bytesMsg := bytes.Join(
		[][]byte{
			pack.SizeToBytes(),
			pack.Bytes(),
		},
		[]byte{},
	)

	skey := encoding.Base64Encode(msg.Body().Hash())
	node.fMapping[skey] = true

	for _, conn := range node.fConnections {
		_, err := conn.Write(bytesMsg)
		if err != nil {
			conn.Close()
		}
	}
}

func (node *sNode) response(nonce []byte, data []byte) {
	node.fMutex.Lock()
	skey := encoding.Base64Encode(nonce)
	ch, ok := node.fActions[skey]
	if !ok {
		return
	}
	node.fMutex.Unlock()

	ch <- data
	close(ch)
}

func (node *sNode) setFunction(name []byte, handle iHandler) INode {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(name)
	node.fHRoutes[skey] = handle
	return node
}

func (node *sNode) getFunction(name []byte) iHandler {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(name)
	f, ok := node.fHRoutes[skey]
	if !ok {
		return nil
	}
	return f
}

func (node *sNode) setAction(nonce []byte) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	node.fActions[skey] = make(chan []byte)
}

func (node *sNode) getAction(nonce []byte) chan []byte {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	ch, ok := node.fActions[skey]
	if !ok {
		panic("undefined key")
	}

	return ch
}

func (node *sNode) delAction(nonce []byte) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	delete(node.fActions, skey)
}

func (node *sNode) setMapping(hash []byte) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	if uint64(len(node.fMapping)) > node.fClient.Settings().Get(settings.SizeMapp) {
		for k := range node.fMapping {
			delete(node.fMapping, k)
			break
		}
	}

	skey := encoding.Base64Encode(hash)
	node.fMapping[skey] = true
}

func (node *sNode) inMapping(hash []byte) bool {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(hash)
	_, ok := node.fMapping[skey]
	return ok
}

func (node *sNode) hasMaxConnSize() bool {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	return uint64(len(node.fConnections)) > node.fClient.Settings().Get(settings.SizeConn)
}

func (node *sNode) setConnection(id string, conn net.Conn) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fConnections[id] = conn
}

func (node *sNode) getConnection(id string) (net.Conn, bool) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	conn, ok := node.fConnections[id]
	return conn, ok
}
