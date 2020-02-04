package gopeer

import (
    "crypto/rsa"
    "errors"
    "net"
    "strings"
    "time"
)

// Create new Listener by address "ipv4:port".
func NewListener(address string) *Listener {
    splited := strings.Split(address, ":")
    if len(splited) != 2 {
        return nil
    }
    return &Listener{
        Address: Address{
            Ipv4: splited[0],
            Port: ":" + splited[1],
        },
        Clients: make(map[string]*Client),
    }
}

// Create new client in listener by private key.
func (listener *Listener) NewClient(private *rsa.PrivateKey) *Client {
    public := &private.PublicKey
    hash := HashPublic(public)
    listener.Clients[hash] = &Client{
        Hashname: hash,
        Keys: Keys{
            Private: private,
            Public:  public,
        },
        Address:     listener.Address.Ipv4 + listener.Address.Port,
        Connections: make(map[string]*Connect),
    }
    return listener.Clients[hash]
}

// Open connection for receiving data.
func (listener *Listener) Open() *Listener {
    var err error
    listener.listen, err = net.Listen("tcp", settings.TEMPLATE+listener.Address.Port)
    if err != nil {
        return nil
    }
    return listener
}

// Run handle server for listening packages.
func (listener *Listener) Run(handleServer func(*Client, *Package)) *Listener {
    go runServer(handleServer, listener)
    return listener
}

// Close listener connection.
func (listener *Listener) Close() {
    if listener == nil {
        return
    }
    listener.listen = nil
}

// Check if user saved in client data.
func (client *Client) InConnections(hash string) bool {
    if _, ok := client.Connections[hash]; ok {
        return true
    }
    return false
}

// Switcher function about GET and SET options.
// GET: accept package and send response;
// SET: accept package;
func (client *Client) HandleAction(title string, pack *Package, handleGet func(*Client, *Package) string, handleSet func(*Client, *Package)) bool {
    switch pack.Head.Title {
    case title:
        switch pack.Head.Option {
        case settings.OPTION_GET:
            data := handleGet(client, pack)
            client.Send(&Package{
                To: To{
                    Receiver: Receiver{
                        Hashname: pack.From.Sender.Hashname,
                    },
                    Address: pack.From.Address,
                },
                Head: Head{
                    Title:  title,
                    Option: settings.OPTION_SET,
                },
                Body: Body{
                    Data: data,
                },
            })
            return true
        case settings.OPTION_SET:
            handleSet(client, pack)
            client.Connections[pack.From.Sender.Hashname].waiting <- true
            return true
        }
    }
    return false
}

// Disconnect from user.
// Send package for disconnect.
// If the user is not responding: delete in local data.
func (client *Client) Disconnect(dest *Destination) error {
    hash := HashPublic(dest.Public)
    _, err := client.SendTo(dest, &Package{
        Head: Head{
            Title:  settings.TITLE_DISCONNECT,
            Option: settings.OPTION_GET,
        },
    })
    delete(client.Connections, hash)
    return err
}

// Connect to user.
// Create local data with parameters.
// After sending GET and receiving SET package, set Connected = true.
func (client *Client) Connect(dest *Destination) error {
    var (
        hash               = HashPublic(dest.Public)
        session            = GenerateRandomBytes(32)
        lastHash           = settings.GENESIS
        prevSession []byte = nil
    )
    if client.InConnections(hash) {
        lastHash = client.Connections[hash].lastHash
        prevSession = client.Connections[hash].Session
    }
    client.Connections[hash] = &Connect{
        connected:   false,
        prevSession: prevSession,
        Session:     session,
        Address:     dest.Address,
        lastHash:    lastHash,
        Public:      dest.Public,
        waiting:     make(chan bool),
    }
    _, err := client.SendTo(dest, &Package{
        Head: Head{
            Title:  settings.TITLE_CONNECT,
            Option: settings.OPTION_GET,
        },
        Body: Body{
            Data: Base64Encode(EncryptRSA(dest.Public, session)),
        },
    })
    return err
}

// Load file from node.
// Input = name file in node side.
// Output = result name file in our side.
func (client *Client) LoadFile(dest *Destination, input string, output string) error {
    hash := HashPublic(dest.Public)
    if !client.InConnections(hash) {
        return errors.New("client not connected")
    }

    client.Connections[hash].transfer.isBlocked  = true
    client.Connections[hash].transfer.outputFile = output

    for id := uint32(0); client.isBlocked(hash); id++ {
        client.SendTo(dest, &Package{
            Head: Head{
                Title: settings.TITLE_FILETRANSFER,
                Option: settings.OPTION_GET,
            },
            Body: Body{
                Data: string(PackJSON(FileTransfer{
                    Head: HeadTransfer{
                        Id: id,
                        Name: input,
                    },
                })),
            },
        })
    }

    nullpack := string(PackJSON(FileTransfer{
        Head: HeadTransfer {
            IsNull: true,
        },
    }))

    client.SendTo(dest, &Package{
        Head: Head{
            Title: settings.TITLE_FILETRANSFER,
            Option: settings.OPTION_GET,
        },
        Body: Body{
            Data: nullpack,
        },
    })
    return nil
}

// Send by Destination struct{
//  Address
//  Public Key
// )
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {
    pack.To.Receiver.Hashname = HashPublic(dest.Public)
    pack.To.Address = dest.Address
    return client.Send(pack)
}

// Send package.
// Check if pack is not null and receive user in saved data.
// Append headers and confirm package.
// Send package.
// If option package is GET, then get response.
// If no response received, then use retrySend() function.
func (client *Client) Send(pack *Package) (*Package, error) {
    if pack == nil {
        return nil, errors.New("pack is null")
    }
    if pack.To.Receiver.Hashname == client.Hashname {
        return nil, errors.New("sender and receiver is one person")
    }
    if !client.InConnections(pack.To.Receiver.Hashname) {
        return nil, errors.New("receiver not in connections")
    }
    if client.isBlocked(pack.To.Receiver.Hashname) && pack.Head.Title != settings.TITLE_FILETRANSFER {
        return nil, errors.New("connections is blocked")
    }
    client.confirmPackage(client.appendHeaders(pack))
    var (
        savedPack = pack
        hash      = pack.To.Receiver.Hashname
    )
    if encPack := client.encryptPackage(pack); encPack != nil {
        pack = encPack
    }
    conn, err := net.Dial("tcp", pack.To.Address)
    if err != nil {
        delete(client.Connections, hash)
        return nil, err
    }
    conn.Write(EncryptAES([]byte(settings.NOISE), PackJSON(pack)))
    conn.Close()
    if savedPack.Head.Option == settings.OPTION_GET {
        select {
        case <-client.Connections[hash].waiting:
            err = nil
        case <-time.After(time.Duration(settings.RETRY_TIME) * time.Second):
            if savedPack.isLasthash() {
                return savedPack, err
            }
            err = client.retrySend(savedPack)
        }
    }
    return savedPack, err
}
