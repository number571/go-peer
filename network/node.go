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
	fClient      local.IClient
	fPReceiver   crypto.IPubKey
	fF2F         iF2F
	fMutex       sync.Mutex
	fListener    net.Listener
	fHRoutes     map[string]iHandler
	fMapping     map[string]bool
	fConnections map[string]net.Conn
	fActions     map[string]chan []byte
}

// Create client by private key as identification.
func NewNode(client local.IClient) INode {
	if client == nil {
		return nil
	}

	pseudo := crypto.NewPrivKey(client.PubKey().Size())
	return &sNode{
		fClient:      client,
		fPReceiver:   pseudo.PubKey(),
		fHRoutes:     make(map[string]iHandler),
		fMapping:     make(map[string]bool),
		fConnections: make(map[string]net.Conn),
		fActions:     make(map[string]chan []byte),
		fF2F: &sF2F{
			fEnabled: false,
			fFriends: make(map[string]crypto.IPubKey),
		},
	}
}

// Close listener and current connections.
func (node *sNode) Close() {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	node.fListener.Close()
	for id, conn := range node.fConnections {
		conn.Close()
		delete(node.fConnections, id)
	}
}

// Return client structure.
func (node *sNode) Client() local.IClient {
	return node.fClient
}

// Return f2f structure.
func (node *sNode) F2F() iF2F {
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
	node.setFunction(title, handle)
	return node
}

// Send message by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (node *sNode) Request(route local.IRoute, msg local.IMessage) ([]byte, error) {
	var (
		result []byte
		err    error
	)

	var (
		waitTime = time.Duration(node.Client().Settings().Get(settings.TimeWait))
		retryNum = node.Client().Settings().Get(settings.SizeRtry)
		counter  = uint64(0)
	)

	routeMsg, session := node.Client().Encrypt(route, msg)
	if routeMsg == nil {
		return nil, errors.New("psender is nil and routes not nil")
	}

	node.setAction(session)
	defer node.delAction(session)

LOOP:
	for counter = 0; counter < retryNum; counter++ {
		node.send(routeMsg)

		select {
		case result = <-node.getAction(session):
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
	var (
		counter  = uint64(0)
		retryNum = node.Client().Settings().Get(settings.SizeRtry)
		conn, _  = node.getConnection(id)
	)

	defer func() {
		node.Disconnect(id)
	}()

	for {
		if counter == retryNum {
			break
		}

		msg := node.readMessage(conn)

	REPEAT:
		if msg == nil {
			counter++
			continue
		}

		counter = 0

		// check message in mapping by hash
		if node.inMapping(msg.Body().Hash()) {
			continue
		}
		node.setMapping(msg.Body().Hash())

		// redirect this message to connections
		node.send(msg)

		// try decrypt message
		decMsg, title := node.Client().Decrypt(msg)
		if decMsg == nil {
			continue
		}

		// if this message is just route message
		// then try procedures again
		routeMsg := node.Client().Settings().Get(settings.MaskRout)

		// if is route package then
		// 1/2 generate new pseudo-package and sleep rand time
		// unpack and send new version of package
		if bytes.Equal(title, encoding.Uint64ToBytes(routeMsg)) {
			rand := crypto.NewPRNG()
			if rand.Uint64()%2 == 0 {
				// send pseudo message
				pMsg, _ := node.Client().Encrypt(
					local.NewRoute(node.fPReceiver, nil, nil),
					local.NewMessage(
						rand.Bytes(16),
						rand.Bytes(calcRandSize(len(decMsg.Body().Data()))),
					),
				)
				node.send(pMsg)
				// sleep random milliseconds
				wtime := node.Client().Settings().Get(settings.TimePsdo)
				time.Sleep(time.Millisecond * calcRandTime(wtime))
			}
			msg = local.LoadPackage(decMsg.Body().Data()).ToMessage()
			goto REPEAT
		}

		// if mode is friend-to-friend and sender not in list of f2f
		// then pass this request
		sender := crypto.LoadPubKey(decMsg.Head().Sender())
		if node.F2F().Status() && !node.F2F().InList(sender) {
			continue
		}

		// send message to handler
		node.handleFunc(decMsg, title)
	}
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
		local.NewRoute(crypto.LoadPubKey(msg.Head().Sender()), nil, nil),
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

func (node *sNode) send(msg local.IMessage) {
	node.fMutex.Lock()
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
	for _, cn := range node.fConnections {
		go cn.Write(bytesMsg)
	}
}

func (node *sNode) response(nonce []byte, data []byte) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	if _, ok := node.fActions[skey]; ok {
		node.fActions[skey] <- data
	}
}

func (node *sNode) setFunction(name []byte, handle iHandler) {
	node.fMutex.Lock()
	defer node.fMutex.Unlock()

	skey := encoding.Base64Encode(name)
	node.fHRoutes[skey] = handle
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
		return make(chan []byte)
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

	if uint64(len(node.fMapping)) > node.Client().Settings().Get(settings.SizeMapp) {
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

	return uint64(len(node.fConnections)) > node.Client().Settings().Get(settings.SizeConn)
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

func calcRandSize(len int) uint64 {
	ulen := uint64(len)
	rand := crypto.NewPRNG()
	return ulen + rand.Uint64()%(10<<10) // +[0;10]KiB
}

func calcRandTime(wtime uint64) time.Duration {
	rand := crypto.NewPRNG()
	return time.Duration(rand.Uint64() % wtime) // +[0;wtime]MS
}
