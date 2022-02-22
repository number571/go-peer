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
	_ Node = &NodeT{}
)

// Basic structure for network use.
type NodeT struct {
	client      local.Client
	preceiver   crypto.PubKey
	f2f         F2F
	mutex       sync.Mutex
	listener    net.Listener
	hroutes     map[string]Handler
	mapping     map[string]bool
	connections map[string]net.Conn
	actions     map[string]chan []byte
}

// Create client by private key as identification.
func NewNode(client local.Client) Node {
	if client == nil {
		return nil
	}

	pseudo := crypto.NewPrivKey(client.PubKey().Size())
	return &NodeT{
		client:      client,
		preceiver:   pseudo.PubKey(),
		hroutes:     make(map[string]Handler),
		mapping:     make(map[string]bool),
		connections: make(map[string]net.Conn),
		actions:     make(map[string]chan []byte),
		f2f: &F2FT{
			friends: make(map[string]crypto.PubKey),
		},
	}
}

// Close listener and current connections.
func (node *NodeT) Close() {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.listener.Close()
	for id, conn := range node.connections {
		conn.Close()
		delete(node.connections, id)
	}
}

// Return client structure.
func (node *NodeT) Client() local.Client {
	return node.client
}

// Return f2f structure.
func (node *NodeT) F2F() F2F {
	return node.f2f
}

// Turn on listener by address.
// Client handle function need be not null.
func (node *NodeT) Listen(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listen.Close()

	node.listener = listen
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
		id := crypto.RandString(rsize)

		node.setConnection(id, conn)
		go node.handleConn(id)
	}

	return nil
}

// Add function to mapping for route use.
func (node *NodeT) Handle(title []byte, handle Handler) Node {
	node.setFunction(title, handle)
	return node
}

// Send message by public key of receiver.
// Function supported multiple routing with pseudo sender.
func (node *NodeT) Request(route local.Route, msg local.Message) (Response, error) {
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
func (node *NodeT) Connections() []string {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	var list []string
	for addr := range node.connections {
		list = append(list, addr)
	}

	return list
}

// Check the existence of an address in the list of connections.
func (node *NodeT) InConnections(address string) bool {
	_, ok := node.getConnection(address)
	return ok
}

// Connect to node by address.
// Client handle function need be not null.
func (node *NodeT) Connect(address string) error {
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
func (node *NodeT) Disconnect(address string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	conn, ok := node.connections[address]
	if ok {
		conn.Close()
	}

	delete(node.connections, address)
}

func (node *NodeT) handleConn(id string) {
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
		if msg == nil {
			counter++
			continue
		}

		counter = 0

	REPEAT:
		// check message in mapping by hash
		if node.inMapping(msg.Body.Hash) {
			continue
		}
		node.setMapping(msg.Body.Hash)

		// redirect this message to connections
		node.send(msg)

		// try decrypt message
		decMsg := node.client.Decrypt(msg)
		if decMsg == nil {
			continue
		}

		// decrypt sender of message
		sender := crypto.LoadPubKey(decMsg.Head.Sender)
		if sender == nil {
			continue
		}

		// if mode is friend-to-friend and sender not in list of f2f
		// then pass this request
		if node.f2f.Status() && !node.f2f.InList(sender) {
			continue
		}

		// if this message is just route message
		// then try procedures again
		title, data := decMsg.Export()
		routeMsg := node.Client().Settings().Get(settings.MaskRout)

		// if is route package then
		// 1/2 generate new pseudo-package and sleep rand time
		// unpack and send new version of package
		if bytes.Equal(title, encoding.Uint64ToBytes(routeMsg)) {
			if crypto.RandUint64()%2 == 0 {
				// send pseudo message
				pMsg, _ := node.Client().Encrypt(
					local.NewRoute(node.preceiver, nil, nil),
					local.NewMessage(
						crypto.RandBytes(16),
						crypto.RandBytes(calcRandSize(len(data))),
					),
				)
				node.send(pMsg)
				// sleep random milliseconds
				wtime := node.Client().Settings().Get(settings.TimePsdo)
				time.Sleep(time.Millisecond * calcRandTime(wtime))
			}
			msg = local.Package(data).Deserialize()
			goto REPEAT
		}

		// send message to handler
		decMsg.Body.Data = data
		node.handleFunc(decMsg, title)
	}
}

func (node *NodeT) handleFunc(msg local.Message, title []byte) {
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
			msg.Body.Data,
		)
		return
	}

	// send response
	f := node.getFunction(title)
	if f == nil {
		return
	}

	rmsg, _ := node.client.Encrypt(
		local.NewRoute(crypto.LoadPubKey(msg.Head.Sender), nil, nil),
		local.NewMessage(
			bytes.Join(
				[][]byte{
					respBytes,
					msg.Head.Session,
					title,
				},
				[]byte{},
			),
			f(node.client, msg),
		),
	)
	node.send(rmsg)
}

func (node *NodeT) send(msg local.Message) {
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

	skey := encoding.Base64Encode(msg.Body.Hash)
	node.mapping[skey] = true
	for _, cn := range node.connections {
		go cn.Write(bytesMsg)
	}
}

func (node *NodeT) response(nonce []byte, data []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	if _, ok := node.actions[skey]; ok {
		node.actions[skey] <- data
	}
}

func (node *NodeT) setFunction(name []byte, handle Handler) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(name)
	node.hroutes[skey] = handle
}

func (node *NodeT) getFunction(name []byte) Handler {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(name)
	f, ok := node.hroutes[skey]
	if !ok {
		return nil
	}
	return f
}

func (node *NodeT) setAction(nonce []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	node.actions[skey] = make(chan []byte)
}

func (node *NodeT) getAction(nonce []byte) chan []byte {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	ch, ok := node.actions[skey]
	if !ok {
		return make(chan []byte)
	}

	return ch
}

func (node *NodeT) delAction(nonce []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(nonce)
	delete(node.actions, skey)
}

func (node *NodeT) setMapping(hash []byte) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	if uint64(len(node.mapping)) > node.client.Settings().Get(settings.SizeMapp) {
		for k := range node.mapping {
			delete(node.mapping, k)
			break
		}
	}

	skey := encoding.Base64Encode(hash)
	node.mapping[skey] = true
}

func (node *NodeT) inMapping(hash []byte) bool {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	skey := encoding.Base64Encode(hash)
	_, ok := node.mapping[skey]
	return ok
}

func (node *NodeT) hasMaxConnSize() bool {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	return uint64(len(node.connections)) > node.client.Settings().Get(settings.SizeConn)
}

func (node *NodeT) setConnection(id string, conn net.Conn) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.connections[id] = conn
}

func (node *NodeT) getConnection(id string) (net.Conn, bool) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	conn, ok := node.connections[id]
	return conn, ok
}

func calcRandSize(len int) uint64 {
	ulen := uint64(len)
	rand := crypto.RandUint64() % (10 << 10)
	return ulen + rand // +[0;10]KiB
}

func calcRandTime(wtime uint64) time.Duration {
	rtime := crypto.RandUint64()
	return time.Duration(rtime % wtime) // +[0;wtime]MS
}
