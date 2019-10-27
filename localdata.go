package gopeer

import (
    "os"
    "net"
    "fmt"
    "time"
    "bufio"
    "bytes"
    "strings"
    "math/rand"
    "crypto/md5"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
)

func runServer(handle func(*Node, *Package), node *Node, conn net.Conn) {
    var sender RelationType

for {
    pack, role := node.readPackage(conn)

redirectPackage:
    // printJsonPackage(pack)

    // Check encrypted package.
    switch {
        case pack == nil: return
        case role == RelationNode && !setting.HAS_FRIENDS &&
             node.InAccessList(AccessDenied, pack.From.Address): return
        case role == RelationNode &&  setting.HAS_FRIENDS && 
            !node.InAccessList(AccessAllowed, pack.From.Address): return
        case node.Setting.ReadOnly == ReadNode && role == RelationHandle: return
        case node.Setting.ReadOnly == ReadHandle && role == RelationNode: return
    }

    // Decrypt package.
    if setting.HAS_CRYPTO {
        node.decryptPackage(pack)
        if packageIsValid(pack) != 0 {
            return
        }
    }

    // Check decrypted package.
    switch {
        case pack == nil: return
        case setting.NETWORK_NAME != pack.Info.NET: return
    }

    // printJsonPackage(pack)

    // Read package.
    switch pack.Head.Title {
        case setting.TITLE_CONNECT:
            switch pack.Head.Mode {
                case setting.MODE_TEST:
                    node.deleteTestConnection(pack.From.Address)
                    if role == RelationNode { return } else { continue }

                case setting.MODE_REMV:
                    node.receiveDisconnect(pack.From.Address)
                    if role == RelationNode { return } else { continue }

                case setting.MODE_READ:
                    if role == RelationHidden {
                        node.receiveHiddenConnect(pack)
                        if sender == RelationNode && !setting.IS_DISTRIB { 
                            return 
                        }
                        continue
                    } else if role == RelationNode {
                        node.receiveNodeConnect(pack)
                        return
                    } else if role == RelationHandle {
                        node.receiveHandleConnect(pack, conn)
                        continue
                    }
                    return
                case setting.MODE_READ_MERG:
                    if setting.IS_DISTRIB {
                        node.receiveMergeConnect(pack, role, conn)
                    }
                    if role == RelationNode { return } else { continue }

                case setting.MODE_SAVE:
                    if role == RelationHidden {
                        node.saveHiddenConnect(pack)
                        if sender == RelationNode && !setting.IS_DISTRIB { 
                            return 
                        }
                        continue
                    } else if role == RelationNode {
                        node.saveNodeConnect(pack)
                        return
                    } else if role == RelationHandle {
                        node.saveHandleConnect(pack, conn)
                        continue
                    }
                    return
                case setting.MODE_SAVE_MERG:
                    if setting.IS_DISTRIB {
                        node.saveMergeConnect(pack, role, conn)
                    }
                    if role == RelationNode { return } else { continue }
            }

        case setting.TITLE_REDIRECT:
            if !setting.HAS_ROUTING { return }
            node.packReceived(pack.From.Address)
            switch pack.Head.Mode {
                case setting.MODE_DISTRIB:
                    if !setting.IS_DISTRIB { return }
                    if new_pack, ok := node.packageSentToMe(pack); ok {
                        pack = new_pack
                        goto redirectPackage
                    }
                    node.SendRedirect(pack)
                    if role == RelationNode { return } else { continue }

                case setting.MODE_DISTRIB_READ:
                    if !setting.IS_DISTRIB { return }
                    node.Send(&Package{
                        To: To{
                            Address: node.AddressByHashname(pack.Body.Desc[0]),
                        },
                        Head: Head{
                            Title: setting.TITLE_REDIRECT,
                            Mode: setting.MODE_DISTRIB_SAVE,
                        },
                        Body: Body{
                            Data: [DATA_SIZE]string{pack.Body.Data[0]},
                        },
                    })
                    if role == RelationNode { return } else { continue }

                case setting.MODE_DISTRIB_SAVE:
                    if !setting.IS_DISTRIB { return }
                    pack = bytesToPack(base64DecodeString(pack.Body.Data[0]))
                    role = RelationHidden
                    goto redirectPackage

                case setting.MODE_DECENTR:
                    if !setting.IS_DECENTR { return }
                    var (
                        packageID = pack.Body.Desc[1]
                        address = pack.Body.Desc[0]
                    )
                    myAddress := node.Address.IPv4 + node.Address.Port
                    if setting.HAS_CRYPTO {
                        myAddress = node.Hashname
                    }
                    if node.inTestConnection(packageID) { return }
                    node.newTestConnection(packageID)
                    go func() {
                        time.Sleep(time.Second * setting.PACK_TIME)
                        node.deleteTestConnection(packageID)
                    }()
                    if node.IsHandle(pack.From.Address) {
                        sender = RelationHandle
                    } else {
                        sender = RelationNode
                    }
                    if address == myAddress {
                        pack = bytesToPack(base64DecodeString(pack.Body.Data[0]))
                        role = RelationHidden
                        goto redirectPackage
                    }
                    node.SendRedirect(pack)
                    if sender == RelationNode { return } else { continue }
            }
    }

    // Check package after read.
    switch {
        case !node.InConnections(pack.From.Address): return
        case role != RelationHidden: node.packReceived(pack.From.Address)
    }

    handle(node, pack)
    if role == RelationNode {
        conn.Close()
        return
    }
}
}

func (node *Node) runInit(handle func(*Node)) *Node {
    handle(node)
    return node
}

func (node *Node) runClient(handle func(*Node, []string)) *Node {
    switch {
        case setting.HAS_CRYPTO && node.Keys.Private == nil:
            return nil
    }
    for {
        handle(node, strings.Split(inputString(), " "))
    }
    return node
}

func (node *Node) runServer(handle func(*Node, *Package), handleInit func(*Node)) *Node {
    switch {
        case setting.HAS_CRYPTO && node.Keys.Private == nil:
            return nil
        case node.Setting.Listen == nil:
            node.Setting.HandleServer = handle
            return node.runInit(handleInit)
    }
    node.runInit(handleInit)
    for {
        conn, err := node.Setting.Listen.Accept()
        if err != nil {
            break
        }
        go runServer(handle, node, conn)
    }
    return node
}

func (node *Node) packReceived(addr string) *Node {
    return node.Send(&Package{
        To: To{
            Address: addr,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_TEST,
        },
    })
}

func (node *Node) readPackage(conn net.Conn) (*Package, RelationType) {
    var (
        pack Package
        message string
        role = RelationNode
        buffer = make([]byte, setting.BUFF_SIZE)
    )
    for {
        length, err := conn.Read(buffer)
        if err != nil || length == 0 { break }
        message += string(buffer[:length])
        if len(message) > setting.MAXSIZE_PACKAGE {
            return nil, role
        }
        if strings.HasSuffix(message, setting.END_BYTES) {
            message = strings.TrimSuffix(message, setting.END_BYTES)
            role = RelationHandle
            break
        }
    }
    err := json.Unmarshal([]byte(message), &pack)
    if err != nil {
        return nil, role
    }
    return &pack, role
}

func (node *Node) findRouting(pack *Package) *Package {
    if pack == nil || node.IsAmI(pack.To.Address) { return nil }
    return &Package{
        Head: Head{
            Title: setting.TITLE_REDIRECT,
            Mode: setting.MODE_DECENTR,
        },
        Body: Body{
            Data: [DATA_SIZE]string{base64.StdEncoding.EncodeToString(packToBytes(node.formatPackage(pack, true)))},
            Desc: [DATA_SIZE]string{
                pack.To.Address,
                base64.StdEncoding.EncodeToString(GenerateRandomBytes(setting.PID_SIZE)),
            },
        },
    }
}

func (node *Node) findInOnionRouting(pack *Package) *Package {
    if pack == nil || node.IsAmI(pack.To.Address) { return nil }
    list := node.GetConnections(RelationHandle)
    return node.onionRouting(&Package{
        To: To{
            Address: list[randomInt(0, len(list))],
        },
        Head: Head{
            Title: setting.TITLE_REDIRECT,
            Mode: setting.MODE_DISTRIB_READ,
        },
        Body: Body{
            Data: [DATA_SIZE]string{base64.StdEncoding.EncodeToString(packToBytes(node.formatPackage(pack, true)))},
            Desc: [DATA_SIZE]string{
                pack.To.Address,
            },
        },
    }, RelationHandle)
}

func (node *Node) onionRouting(pack *Package, role RelationType) *Package {
    if pack == nil || node.IsAmI(pack.To.Address) ||
        !node.InConnections(pack.To.Address) { return nil }

    var (
        newPack string
        from string
    )

    list := node.randomRouting(pack.To.Address, role)
    prevPack := node.wrapPackage(pack)

    if setting.HAS_CRYPTO {
        newPack = base64.StdEncoding.EncodeToString(EncryptAES(node.Network.Connections[pack.To.Address].Session, prevPack))
        from = base64.StdEncoding.EncodeToString(EncryptRSA(
            node.Network.Connections[pack.To.Address].Public, 
            []byte(node.Address.IPv4 + node.Address.Port),
        ))
    } else {
        newPack = base64.StdEncoding.EncodeToString(prevPack)
        from = node.Address.IPv4 + node.Address.Port
    }

    return &Package{
        Head: Head{
            Title: setting.TITLE_REDIRECT,
            Mode: setting.MODE_DISTRIB,
        },
        Body: Body{
            Data: [DATA_SIZE]string{newPack},
            Desc: [DATA_SIZE]string{
                strings.Join(list, setting.SEPARATOR),
                from,
            },
        },
    }
}

func (node *Node) onionPackage(pack *Package) *Package {
    newList := strings.Split(pack.Body.Desc[0], setting.SEPARATOR)
    if len(newList) == 0 {
        return nil
    }
    if setting.HAS_CRYPTO {
        pack.To.Address = string(node.DecryptRSA(base64DecodeString(newList[0])))
    } else {
        pack.To.Address = newList[0]
    }
    pack.Body.Desc[0] = strings.Join(newList[1:], setting.SEPARATOR)
    return pack
}

func randomInt(min, max int) int {
    return rand.Intn(max - min) + min
}

func (node *Node) randomRouting(receiver string, role RelationType) []string {
    var (
        list []string
        index int
    )

    for addr, conn := range node.Network.Connections {
        if addr == receiver { continue }
        if conn.Relation != role { continue }
        if index == setting.ROUTE_NUM { break }
        list = append(list, addr)
        index++
    }

    list = append(shuffle(list), receiver)

    if setting.HAS_CRYPTO {
        var (
            length = len(list)
            newList = make([]string, length)
        )
        newList[0] = base64.StdEncoding.EncodeToString(EncryptRSA(node.Keys.Public, []byte(list[0])))
        for i := 0; i < length - 1; i++ {
            newList[i+1] = base64.StdEncoding.EncodeToString(EncryptRSA(node.Network.Connections[list[i]].Public, []byte(list[i+1])))
        }
        list = newList
    }

    return list
}

func shuffle(slice []string) []string {
    length := len(slice)
    mod := uint64(length)
    randomInts := GenerateRandomIntegers(length)
    for i := 0; i < length; i++ {
        j := randomInts[i] % mod
        slice[i], slice[j] = slice[j], slice[i]
    }
    return slice
}

func bytesToPack(data []byte) *Package {
    var pack Package
    err := json.Unmarshal(data, &pack)
    if err != nil {
        return nil
    }
    return &pack
}

func (node *Node) packageSentToMe(pack *Package) (*Package, bool) {
    var (
        newPack []byte
    )
    if setting.HAS_CRYPTO {
        addr := node.DecryptRSA(base64DecodeString(pack.Body.Desc[1]))
        if addr == nil {
            return nil, false
        }
        newPack = DecryptAES(node.Network.Connections[string(addr)].Session, base64DecodeString(pack.Body.Data[0]))
        if newPack == nil {
            return nil, false
        }
    } else {
        // If list not null then redirect package to another node.
        if len(pack.Body.Desc[0]) != 0 {
            return nil, false
        }
        newPack = base64DecodeString(pack.Body.Data[0])
    }
    return bytesToPack(newPack), true
}

func (node *Node) connectToNode(addr string) *Node {
    node.Network.Connections[addr] = &Connect{
        Relation: RelationNode,
    }
    return node.Send(&Package{
        To: To{
            Address: addr,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_READ,
        },
    })
}

func (node *Node) connectToFriend(addr, session string) *Node {
    if setting.HAS_FRIENDS {
        node.AppendToAccessList(AccessAllowed, addr)
    }
    node.Network.Connections[addr] = &Connect{
        Relation: RelationNode,
        Session: HashSum([]byte(session)),
    }
    return node
}

func (node *Node) handleConnect(handle func(*Node, *Package), addr string) *Node {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil
    }
    node.Network.Connections[addr] = &Connect{
        Relation: RelationHandle,
        Link: conn,
    }
    go runServer(handle, node, conn)
    return node.Send(&Package{
        To: To{
            Address: addr,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_READ,
        },
    })
}

func (node *Node) formatPackage(pack *Package, hidden bool) *Package {
    pack.Info.NET = setting.NETWORK_NAME

    if setting.HAS_CRYPTO && hidden {
        pack.From.Address = node.Hashname
    } else {
        pack.From.Address = node.Address.IPv4 + node.Address.Port
    }

    if setting.HAS_CRYPTO {
        pack.From.Hashname = node.Hashname
        pack.From.Public = base64.StdEncoding.EncodeToString([]byte(node.StringPublic()))

        pack.Body.Time = time.Now().Format(time.RFC1123)
        pack.Body.Hash = ""
        pack.Body.Sign = ""

        hash_pack := hashPack(pack)
        pack.Body.Hash = base64.StdEncoding.EncodeToString(hash_pack)
        pack.Body.Sign = base64.StdEncoding.EncodeToString(node.Sign(hash_pack))

        // printJsonPackage(pack)
        node.encryptPackage(pack)
    }

    return pack
}

func (node *Node) wrapPackage(pack *Package) []byte {
    return packToBytes(node.formatPackage(pack, false))
}

func (node *Node) setTimePackage(pack *Package, titleIsConnect bool) {
    if !titleIsConnect {
        go func() {
            node.newTestConnection(pack.To.Address)
            time.Sleep(time.Second * setting.WAIT_TIME)
            if _, ok := node.Setting.TestConnections[pack.To.Address]; ok {
                node.deleteTestConnection(pack.To.Address)
                node.receiveDisconnect(pack.To.Address)
            }
        }()
    }
}

func (node *Node) sendToNode(pack *Package) *Node {
    conn, err := net.Dial("tcp", pack.To.Address)
    if err != nil {
        if !node.InConnections(pack.To.Address) {
            return nil
        }
        if node.Network.Connections[pack.To.Address].Relation == RelationHidden { 
            return node
        }
        node.receiveDisconnect(pack.To.Address)
        return nil
    }
    defer conn.Close()
    titleIsConnect := pack.Head.Title == setting.TITLE_CONNECT
    data := node.wrapPackage(pack)
    conn.Write(data)
    node.setTimePackage(pack, titleIsConnect)
    return node
}

func (node *Node) sendHandle(pack *Package) *Node {
    conn := node.Network.Connections[pack.To.Address].Link
    titleIsConnect := pack.Head.Title == setting.TITLE_CONNECT
    data := node.wrapPackage(pack)
    _, err := conn.Write(bytes.Join(
        [][]byte{data, []byte(setting.END_BYTES)},
        []byte{},
    ))
    if err != nil {
        node.receiveDisconnect(pack.To.Address)
        return nil
    }
    node.setTimePackage(pack, titleIsConnect)
    return node
}

func printJsonPackage(pack *Package) {
    jsonData, _ := json.MarshalIndent(pack, "", "\t")
    fmt.Println(string(jsonData))
}

func (node *Node) encryptPackage(pack *Package) {
    var sessionKey []byte
    if pack.Head.Title == setting.TITLE_CONNECT && (pack.Head.Mode == setting.MODE_READ || pack.Head.Mode == setting.MODE_READ_MERG) {
        return
    } else if pack.Head.Title == setting.TITLE_CONNECT && (pack.Head.Mode == setting.MODE_SAVE || pack.Head.Mode == setting.MODE_SAVE_MERG) {
        sessionKey = []byte(pack.Body.Desc[0])
        pack.Body.Desc[0] = base64.StdEncoding.EncodeToString(EncryptRSA(node.Network.Connections[pack.To.Address].Public, sessionKey))
    } else {
        if !node.InConnections(pack.To.Address) {
            return
        }
        sessionKey = node.Network.Connections[pack.To.Address].Session
        pack.Head.Title = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Head.Title)))
        pack.Head.Mode  = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Head.Mode)))
        pack.Body.Desc[0] = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Desc[0])))
    }
    pack.Body.Data[0] = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Data[0])))
    for i := 1; i < DATA_SIZE; i++ {
        pack.Body.Desc[i] = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Desc[i])))
        pack.Body.Data[i] = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Data[i])))
    }
    pack.Info.NET = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Info.NET)))
    pack.Body.Time = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Time)))
    pack.Body.Hash = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Hash)))
    pack.Body.Sign = base64.StdEncoding.EncodeToString(EncryptAES(sessionKey, []byte(pack.Body.Sign)))
}

func (node *Node) decryptPackage(pack *Package) {
    var sessionKey []byte
    if pack.Head.Title == setting.TITLE_CONNECT && (pack.Head.Mode == setting.MODE_READ || pack.Head.Mode == setting.MODE_READ_MERG) {
        return
    } else if pack.Head.Title == setting.TITLE_CONNECT && (pack.Head.Mode == setting.MODE_SAVE || pack.Head.Mode == setting.MODE_SAVE_MERG) {
        sessionKey = node.DecryptRSA(base64DecodeString(pack.Body.Desc[0]))
        pack.Body.Desc[0] = string(sessionKey)
    } else {
        if !node.InConnections(pack.From.Address) {
            return
        }
        sessionKey = node.Network.Connections[pack.From.Address].Session
        pack.Head.Title = string(DecryptAES(sessionKey, base64DecodeString(pack.Head.Title)))
        pack.Head.Mode  = string(DecryptAES(sessionKey, base64DecodeString(pack.Head.Mode)))
        pack.Body.Desc[0] = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Desc[0])))
    }
    pack.Body.Data[0] = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Data[0])))
    for i := 1; i < DATA_SIZE; i++ {
        pack.Body.Desc[i] = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Desc[i])))
        pack.Body.Data[i] = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Data[i])))
    }
    pack.Info.NET = string(DecryptAES(sessionKey, base64DecodeString(pack.Info.NET)))
    pack.Body.Time = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Time)))
    pack.Body.Hash = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Hash)))
    pack.Body.Sign = string(DecryptAES(sessionKey, base64DecodeString(pack.Body.Sign)))
}

func (node *Node) receiveDisconnect(addr string) {
    if !node.InConnections(addr) { return }
    node.Setting.Mutex.Lock()
    delete(node.Network.Addresses, node.Network.Connections[addr].Hashname)
    delete(node.Network.Connections, addr)
    node.Setting.Mutex.Unlock()
}

func (node *Node) receiveMergeConnect(pack *Package, role RelationType, conn net.Conn) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Hashname] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: role,
        Hashname: pack.From.Hashname, 
        Session: GenerateRandomBytes(setting.SESSION_SIZE),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    if role == RelationHandle {
        node.Network.Connections[pack.From.Address].Link = conn
    }
    node.Setting.Mutex.Unlock()
    node.ConnectToList(strings.Split(pack.Body.Desc[0], setting.SEPARATOR))
    node.Send(&Package{
        To: To{
            Address: pack.From.Address,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_SAVE_MERG,
        },
        Body: Body{
            Desc: [DATA_SIZE]string{
                string(node.Network.Connections[pack.From.Address].Session), 
                strings.Join(node.GetConnections(RelationNode), setting.SEPARATOR),
            },
        },
    })
}

func (node *Node) receiveHiddenConnect(pack *Package) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Address] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: RelationHidden,
        Hashname: pack.From.Address, 
        Session: GenerateRandomBytes(setting.SESSION_SIZE),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    node.Setting.Mutex.Unlock()
    node.SendInitRedirect(&Package{
        To: To{
            Address: pack.From.Address,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_SAVE,
        },
        Body: Body{
            Desc: [DATA_SIZE]string{
                string(node.Network.Connections[pack.From.Address].Session),
            },
        },
    })
}

func (node *Node) receiveNodeConnect(pack *Package) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Hashname] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: RelationNode,
        Hashname: pack.From.Hashname, 
        Session: GenerateRandomBytes(setting.SESSION_SIZE),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    node.Setting.Mutex.Unlock()
    node.Send(&Package{
        To: To{
            Address: pack.From.Address,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_SAVE,
        },
        Body: Body{
            Desc: [DATA_SIZE]string{string(node.Network.Connections[pack.From.Address].Session)},
        },
    })
}

func (node *Node) receiveHandleConnect(pack *Package, conn net.Conn) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Hashname] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: RelationHandle,
        Hashname: pack.From.Hashname, 
        Session: GenerateRandomBytes(setting.SESSION_SIZE),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
        Link: conn,
    }
    node.Setting.Mutex.Unlock()
    node.Send(&Package{
        To: To{
            Address: pack.From.Address,
        },
        Head: Head{
            Title: setting.TITLE_CONNECT,
            Mode: setting.MODE_SAVE,
        },
        Body: Body{
            Desc: [DATA_SIZE]string{string(node.Network.Connections[pack.From.Address].Session)},
        },
    })
}

func (node *Node) saveMergeConnect(pack *Package, role RelationType, conn net.Conn) {
    node.ConnectToList(strings.Split(pack.Body.Desc[1], setting.SEPARATOR))
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Hashname] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: role,
        Hashname: pack.From.Hashname, 
        Session: []byte(pack.Body.Desc[0]),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    if role == RelationHandle {
        node.Network.Connections[pack.From.Address].Link = conn
    }
    node.Setting.Mutex.Unlock()
}

func (node *Node) saveNodeConnect(pack *Package) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Hashname] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: RelationNode,
        Hashname: pack.From.Hashname, 
        Session: []byte(pack.Body.Desc[0]),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    node.Setting.Mutex.Unlock()
}

func (node *Node) saveHandleConnect(pack *Package, conn net.Conn) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Hashname] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: RelationHandle,
        Hashname: pack.From.Hashname, 
        Session: []byte(pack.Body.Desc[0]),
        Link: conn,
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    node.Setting.Mutex.Unlock()
}

func (node *Node) saveHiddenConnect(pack *Package) {
    node.Setting.Mutex.Lock()
    node.Network.Addresses[pack.From.Address] = pack.From.Address
    node.Network.Connections[pack.From.Address] = &Connect{
        Relation: RelationHidden,
        Hashname: pack.From.Address, 
        Session: []byte(pack.Body.Desc[0]),
        Public: ParsePublic(string(base64DecodeString(pack.From.Public))),
    }
    node.Setting.Mutex.Unlock()
}

func (node *Node) inTestConnection(addr string) bool {
    node.Setting.Mutex.Lock()
    _, ok := node.Setting.TestConnections[addr]
    node.Setting.Mutex.Unlock()
    return ok
}

func (node *Node) newTestConnection(addr string) {
    node.Setting.Mutex.Lock()
    node.Setting.TestConnections[addr] = true
    node.Setting.Mutex.Unlock()
}

func (node *Node) deleteTestConnection(addr string) {
    node.Setting.Mutex.Lock()
    delete(node.Setting.TestConnections, addr)
    node.Setting.Mutex.Unlock()
}

func paddingPKCS5(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext) % blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func unpaddingPKCS5(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    if length < unpadding {
        return nil
    }
    return origData[:(length - unpadding)]
}

func md5HashName(data string) string {
    hash := md5.Sum([]byte(data))
    return base64.StdEncoding.EncodeToString(hash[:])
}

func base64DecodeString(data string) []byte {
    result, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return nil
    }
    return result
}

func inputString() string {
    msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    return strings.Replace(msg, "\n", "", -1)
}

func packToBytes(pack *Package) []byte {
    data, err := json.Marshal(pack)
    if err != nil {
        return nil
    }
    return data
}

func hashPack(pack *Package) []byte {
    return HashSum(packToBytes(pack))
}

func sumSHA256(data []byte) []byte {
    hash := sha256.Sum256(data)
    return hash[:]
}

func packageIsValid(pack *Package) uint8 {
    tempPack := *pack

    tempPack.Body.Hash = ""
    tempPack.Body.Sign = ""

    if len(pack.From.Hashname) != LEN_BASE64_SHA256 {
        return 1
    }

    hash := base64.StdEncoding.EncodeToString(hashPack(&tempPack))
    if hash != pack.Body.Hash {
        return 2
    }

    publicString := string(base64DecodeString(pack.From.Public))
    if md5HashName(publicString) != pack.From.Hashname {
        return 3
    }

    public := ParsePublic(publicString)
    if public == nil {
        return 4
    }

    verify := Verify(public, base64DecodeString(pack.Body.Hash), base64DecodeString(pack.Body.Sign))
    if verify != nil {
        return 5
    }

    return 0
}
