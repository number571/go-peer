package gopeer

import (
    "net"
    "sync"
    "strings"
    "encoding/base64"
)

// Create node. Set address and default options.
func NewNode(addr string) *Node {
    if setting.IS_NOTHING { return nil }
    var node = &Node{
        Setting: Setting{
            Mutex: new(sync.Mutex),
            TestConnections: make(map[string]bool),
        },
        Network: Network{
            AccessList: make(map[string]AccessType),
            Connections: make(map[string]*Connect),
        },
    }

    if len(addr) > setting.MAXSIZE_ADDRESS { return nil }
    if addr == setting.CLIENT_NAME {
        node.Address.IPv4 = setting.CLIENT_NAME
        node.Address.Port = base64.StdEncoding.EncodeToString(GenerateRandomBytes(setting.CLIENT_NAME_SIZE))
        return node
    }

    splited := strings.Split(addr, ":")
    if len(splited) != 2 { return nil }

    node.Address.IPv4 = splited[0]
    node.Address.Port = ":" + splited[1]
    return node
}

// Turn server in listening mode.
func (node *Node) Open() *Node {
    var err error
    node.Setting.Listen, err = net.Listen("tcp", setting.TEMPLATE + node.Address.Port)
    if err != nil {
        return nil
    }
    return node
}

// Run server and client applications in parallel.
func (node *Node) Run(handleInit func(*Node), handleServer func(*Node, *Package), handleClient func(*Node, []string)) *Node {
    go node.runServer(handleServer, handleInit)
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
func (node *Node) IsAmI(addr string) bool {
    return node.Address.IPv4 + node.Address.Port == addr
}

// Append addresses to access list.
func (node *Node) AppendToAccessList(access AccessType, addresses ...string) {
    for _, addr := range addresses {
        if node.IsAmI(addr) { continue }
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
    if pack == nil || !setting.HAS_ROUTING { return nil }
    if setting.IS_DECENTR {
        return node.findRouting(pack)
    } else {
        if node.Setting.Listen == nil {
            return node.findInOnionRouting(pack)
        } else {
            return node.onionRouting(pack, RelationNode)
        }
    }
}

// Send package through intermediaries.
func (node *Node) SendRedirect(pack *Package) *Node {
    if pack == nil || !setting.HAS_ROUTING { return nil }
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
    if pack == nil { return nil }
    newPack := *pack
    for addr := range node.Network.Connections {
        if addr == sender { continue }
        newPack.To.Address = addr
        node.Send(&newPack)
    }
    return node
}

// Use function Send for all nodes.
func (node *Node) SendToAll(pack *Package) *Node {
    if pack == nil { return nil }
    newPack := *pack
    for addr := range node.Network.Connections {
        newPack.To.Address = addr
        node.Send(&newPack)
    }
    return node
}

// Send package to another node.
// In crypto mode function uses hashnames, public keys, 
// hashes, timestamp, signatures and encryption.
func (node *Node) Send(pack *Package) *Node {
    if pack == nil || node.IsAmI(pack.To.Address) {
        return nil
    }
    newPack := *pack
    if node.IsHandle(pack.To.Address) {
        return node.sendHandle(&newPack)
    } else {
        return node.sendToNode(&newPack)
    }
    return nil
}

// Connect to node, his nodes and send him connections.
func (node *Node) MergeConnect(addr string) *Node {
    if !setting.IS_DISTRIB || node.IsAmI(addr) {
        return nil
    }
    if node.Setting.Listen == nil {
        conn, err := net.Dial("tcp", addr)
        if err != nil {
            return nil
        }
        node.Network.Connections[addr] = &Connect{
            Relation: RelationHandle,
            Link: conn,
        }
        go runServer(node.Setting.HandleServer, node, conn)
    }
    return node.Send(&Package{
        To: To{
            Address: addr,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_READ_MERG,
        },
        Body: Body{
            Desc: [DATA_SIZE]string{
                strings.Join(node.GetConnections(RelationNode), setting.SEPARATOR),
            },
        },
    })
}

// Connect to hidden friends.
func (node *Node) HiddenConnect(addr string) *Node {
    node.Network.Connections[addr] = &Connect{
        Relation: RelationHidden,
    }
    return node.SendInitRedirect(&Package{
        To: To{
            Address: addr,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_READ,
        },
    })
}

// Connect to many nodes.
func (node *Node) ConnectToList(list ...interface{}) *Node {
    for _, value := range list {
        switch value.(type) {
            case string:
                node.Connect(value.(string))
            case []string:
                data := value.([]string)
                for _, addr := range data {
                    node.Connect(addr)
                }
            case [2]string:
                data := value.([2]string)
                if len(data) != 2 { continue }
                node.Connect(data[0], data[1])
        }
    }
    return node
}

// Connect to node by address.
func (node *Node) Connect(data ...string) *Node {
    // data[0] = address
    // data[1] = password
    if node.IsAmI(data[0]) { return nil }
    switch len(data) {
        case 1:
            if node.Setting.Listen == nil {
                return node.handleConnect(node.Setting.HandleServer, data[0])
            } else {
                return node.connectToNode(data[0])
            }
        case 2:
            if setting.IS_DECENTR {
                return node.connectToFriend(data[0], data[1])
            }
    }
    return nil
}

// Disconnect from node by address.
func (node *Node) Disconnect(addresses ...string) *Node {
    for _, addr := range addresses {
        node.Send(&Package{
            To: To{
                Address: addr,
            },
            Head: Head{
                Title: setting.TITLE_CONNECT,
                Mode: setting.MODE_REMV,
            },
        }).receiveDisconnect(addr)
    }
    return node
}

// Get address by hashname.
func (node *Node) AddressByHashname(hashname string) string {
    for addr, conn := range node.Network.Connections {
        if conn.Hashname == hashname {
            return addr 
        }
    }
    return hashname
}

// Get connected addresses.
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

// Turn off server.
func (node *Node) Close() *Node {
    if node.Setting.Listen == nil {
        return nil
    }
    node.Setting.Listen.Close()
    return node
}
