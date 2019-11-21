package gopeer

import (
	"encoding/base64"
	"net"
	"strings"
	"sync"
)

// Create node. Set address and default options.
func NewNode(addr string) *Node {
	if setting.IS_NOTHING {
		return nil
	}
	var node = &Node{
		Setting: Setting{
			Mutex:           new(sync.Mutex),
			TestConnections: make(map[string]bool),
		},
		Network: Network{
			Addresses:   make(map[string]string),
			AccessList:  make(map[string]AccessType),
			Connections: make(map[string]*Connect),
		},
	}

	if len(addr) > setting.MAXSIZE_ADDRESS {
		return nil
	}
	if addr == setting.CLIENT_NAME {
		node.Address.IPv4 = setting.CLIENT_NAME
		node.Address.Port = base64.StdEncoding.EncodeToString(GenerateRandomBytes(setting.CLIENT_NAME_SIZE))
		return node
	}

	splited := strings.Split(addr, ":")
	if len(splited) != 2 {
		return nil
	}

	node.Address.IPv4 = splited[0]
	node.Address.Port = ":" + splited[1]
	return node
}

// Turn server in listening mode.
func (node *Node) Open() *Node {
	var err error
	node.Setting.Listen, err = net.Listen("tcp", setting.TEMPLATE+node.Address.Port)
	if err != nil {
		return nil
	}
	return node
}

// Run server and client applications in parallel.
func (node *Node) Run(handleServer func(*Node, *Package), handleClient func(*Node)) *Node {
	switch {
	case setting.HAS_CRYPTO && node.Keys.Private == nil:
		return nil
	}
	ch := make(chan bool)
	go node.runServer(handleServer, ch)
	<-ch
	return node.runClient(handleClient)
}

// Set receive package from all users.
func (node *Node) ReadOnly(types ReadonlyType) *Node {
	node.Setting.ReadOnly = types
	return node
}

// Check if address is client.
func (node *Node) IsHidden(addr string) bool {
	return node.InConnections(addr) && node.Network.Connections[addr].Relation == RelationHidden
}

// Check if address is client.
func (node *Node) IsHandle(addr string) bool {
	return node.InConnections(addr) && node.Network.Connections[addr].Relation == RelationHandle
}

// Check if address is node.
func (node *Node) IsNode(addr string) bool {
	return node.InConnections(addr) && node.Network.Connections[addr].Relation == RelationNode
}

// Check if address is my.
func (node *Node) IsMyAddress(addr string) bool {
	return node.Address.IPv4+node.Address.Port == addr
}

// Check if hashname is my.
func (node *Node) IsMyHashname(hashname string) bool {
	return node.Hashname == hashname
}

// Return true if server is not running.
func (node *Node) IAmClient() bool {
	return node.Setting.Listen == nil
}

// Return true if server is running.
func (node *Node) IAmNode() bool {
	return node.Setting.Listen != nil
}

// Append addresses to access list.
func (node *Node) AppendToAccessList(access AccessType, addresses ...string) {
	for _, addr := range addresses {
		if node.IsMyAddress(addr) {
			continue
		}
		node.Network.AccessList[addr] = access
	}
}

// Delete addresses from access list.
func (node *Node) DeleteFromAccessList(addresses ...string) *Node {
	for _, addr := range addresses {
		delete(node.Network.AccessList, addr)
	}
	return node
}

// Check if address exist in access list.
func (node *Node) InAccessList(access AccessType, addr string) bool {
	_, ok := node.Network.AccessList[addr]
	return ok && node.Network.AccessList[addr] == access
}

// Check if address exist in current connections.
func (node *Node) InConnections(addr string) bool {
	_, ok := node.Network.Connections[addr]
	return ok
}

// Create package for send redirect.
func (node *Node) CreateRedirect(pack *Package) *Package {
	if pack == nil || !setting.HAS_ROUTING {
		return nil
	}
	if setting.IS_DECENTR {
		return node.findRouting(pack)
	} else {
		if node.IAmClient() {
			return node.findInOnionRouting(pack)
		} else {
			return node.onionRouting(pack, RelationNode)
		}
	}
}

// Send package through intermediaries.
func (node *Node) SendRedirect(pack *Package) *Node {
	if pack == nil || !setting.HAS_ROUTING {
		return nil
	}
	if setting.IS_DECENTR {
		return node.SendToAllWithout(pack, pack.From.Address)
	} else {
		return node.Send(node.onionPackage(pack))
	}
}

// Create redirect package and send.
func (node *Node) SendInitRedirect(pack *Package) *Node {
	return node.SendRedirect(node.CreateRedirect(pack))
}

// Send package to all connections without sender.
func (node *Node) SendToAllWithout(pack *Package, sender string) *Node {
	if pack == nil {
		return nil
	}
	newPack := *pack
	for addr := range node.Network.Connections {
		if addr == sender {
			continue
		}
		newPack.To.Address = addr
		node.Send(&newPack)
	}
	return node
}

// Use function Send for all nodes.
func (node *Node) SendToAll(pack *Package) *Node {
	return node.SendTo(pack, node.GetConnections(RelationAll)...)
}

// Send package to list of addresses.
func (node *Node) SendTo(pack *Package, addresses ...string) *Node {
	if pack == nil {
		return nil
	}
	newPack := *pack
	for _, addr := range addresses {
		if !node.InConnections(addr) {
			continue
		}
		newPack.To.Address = addr
		node.Send(&newPack)
	}
	return node
}

// Send package to another node.
// In crypto mode function uses hashnames, public keys,
// hashes, timestamp, signatures and encryption.
func (node *Node) Send(pack *Package) *Node {
	if pack == nil || node.IsMyAddress(pack.To.Address) {
		return nil
	}
	newPack := *pack
	if node.IsHandle(pack.To.Address) {
		return node.sendHandle(&newPack)
	}
	return node.sendToNode(&newPack)
}

// Connect to node, his nodes and send to him my connections.
func (node *Node) MergeConnect(addr string) *Node {
	if (!setting.HANDLE_ROUTING && node.IAmClient()) ||
		setting.IS_DECENTR || node.IsMyAddress(addr) {
		return nil
	}
	relation := RelationNode
	if node.IAmClient() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil
		}
		relation = RelationHandle
		node.Network.Connections[addr] = &Connect{
			Relation: RelationHandle,
			Link:     conn,
		}
		go server(node.Setting.HandleServer, node, conn)
	}
	return node.Send(&Package{
		To: To{
			Address: addr,
		},
		Head: Head{
			Title: setting.TITLE_CONNECT,
			Mode:  setting.MODE_READ_MERG,
		},
		Body: Body{
			Desc: [DATA_SIZE]string{
				strings.Join(node.GetConnections(relation), setting.SEPARATOR),
			},
		},
	})
}

// Connect to hidden friends.
func (node *Node) HiddenConnect(addr string) *Node {
	if setting.IS_DISTRIB && node.IAmNode() {
		return nil
	}
	if (!setting.HAS_CRYPTO && node.IsMyAddress(addr)) ||
		(setting.HAS_CRYPTO && node.IsMyHashname(addr)) {
		return nil
	}
	node.Network.Connections[addr] = &Connect{
		Relation: RelationHidden,
	}
	return node.SendInitRedirect(&Package{
		To: To{
			Address: addr,
		},
		Head: Head{
			Title: setting.TITLE_CONNECT,
			Mode:  setting.MODE_READ,
		},
	})
}

// Connect to many nodes.
// If type data = string then { connect(data) }
// If type data = []string then { for d in data: connect(d) }
// If type data = [2]string then { connect(data[0], data[1]) }
// If type data = map[string]string then { for a,p in data: connect(a,p) }
func (node *Node) ConnectToList(data ...interface{}) *Node {
	for _, d := range data {
		switch d.(type) {
		case string:
			node.Connect(d.(string))
		case []string:
			val := d.([]string)
			for _, addr := range val {
				node.Connect(addr)
			}
		case [2]string:
			val := d.([2]string)
			if len(val) != 2 {
				continue
			}
			node.Connect(val[0], val[1])
		case map[string]string:
			val := d.(map[string]string)
			for addr, pasw := range val {
				node.Connect(addr, pasw)
			}
		}
	}
	return node
}

// Connect to node by address and/or password.
// If length data = 2 then connect to friend by password.
func (node *Node) Connect(data ...string) *Node {
	// data[0] = address
	// data[1] = password
	if node.IsMyAddress(data[0]) {
		return nil
	}
	switch len(data) {
	case 1:
		if node.IAmClient() {
			return node.handleConnect(node.Setting.HandleServer, data[0])
		} else {
			return node.connectToNode(data[0])
		}
	case 2:
		return node.connectToFriend(data[0], data[1])
	}
	return nil
}

// Disconnect from node by address.
// Send package to node for delete on his side and
// delete connections on my side.
func (node *Node) Disconnect(addresses ...string) *Node {
	for _, addr := range addresses {
		node.Send(&Package{
			To: To{
				Address: addr,
			},
			Head: Head{
				Title: setting.TITLE_CONNECT,
				Mode:  setting.MODE_REMV,
			},
		})
		node.deleteFromConnection(addr)
	}
	return node
}

// Get connected addresses by relation.
// If relation = RelationAll then function return all connections.
func (node *Node) GetConnections(relation RelationType) []string {
	var list []string
	if relation == RelationAll {
		for address := range node.Network.Connections {
			list = append(list, address)
		}
		return list
	}
	for address, conn := range node.Network.Connections {
		if conn.Relation == relation {
			list = append(list, address)
		}
	}
	return list
}

// Turn off server listening.
func (node *Node) Close() *Node {
	if node == nil {
		return nil
	}
	if node.Setting.Listen == nil {
		return nil
	}
	node.Setting.Listen.Close()
	return node
}
