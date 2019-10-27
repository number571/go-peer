# gopeer
> Framework for create centralized, decentralized, distributed and hybrid networks.

### Template:
```go
package main

import (
    "github.com/number571/gopeer"
)

// Customize the framework settings before start use.
func init() {
    gopeer.SettingsSet(gopeer.SettingsType{
        
    })
}

// Create node in network with address 'IPv4:Port'.
// Open server for listening connections.
// Run server and client parts. handleInit runs after run handleServer.
// Close listening connections after use.
func main() {
    gopeer.NewNode("IPv4:Port").Open().Run(handleInit, handleServer, handleClient).Close()
}

// Actions before using the client interface.
func handleInit(node *gopeer.Node) {
    
}

// Read package from network.
func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {

    }
}

// Read input message. 
// CLI interface. If another interface is used then this function can be void.
func handleClient(node *gopeer.Node, message []string) {
    switch message[0] {
        
    }
}
```
***
### Main settings:
```go
{
    IS_NOTHING: true,
    IS_DISTRIB: false,
    IS_DECENTR: false,
    HAS_CRYPTO: false,
    HAS_ROUTING: false,
    HAS_FRIENDS: false,
    CRYPTO_SPEED: false,
}
```
#### The choice between decentralized and distributed networks. 
* Mode = "IS_NOTHING" (Node can't be created);
* Mode = "IS_DISTRIB";
* Mode = "IS_DECENTR";
  
#### Based on previous networks implemented centralized and hybrid.
* Need use ReadOnly(RelationHandles) for create centralized networks;
* Need use ReadOnly(RelationNodes) for create decentralized/distributed networks;
* By default is hybrid network.
    
#### Can be used cryptography: 
* Mode = "HAS_CRYPTO" (Key exchange, data encryption, digital signatures, packet verification, hash names);
* Symmetric encryption = AES-CBC (The key depends from "CRYPTO_SPEED" option);
* Asymmetric encryption = RSA-OAEP (The key depends from user node);
* Digital signature = RSASSA-PSS;
* Hash function = SHA256;
  
#### Cryptography can be used in fast mode (not recommended). 
* Mode = "CRYPTO_SPEED";
* The AES key becomes 128 (Default: 256) bits;
* The HashSum function becomes H(H(X)) (Default: (H(H(x), x)));
  
#### Can be used routing:
* Mode = "HAS_ROUTING";
* In decentralized network is a broadcast data.
* In distributed network is an onion routing.
  
#### Nodes have friends:
* Mode = "HAS_FRIENDS";
* In this mode nodes can be connected only by passwords.
***
### Constants:
```go
const(
    LEN_BASE64_SHA256 = 24
    DATA_SIZE = 3
)
```
* LEN_BASE64_SHA256 = len(base64(sha256(X)));
* DATA_SIZE = length arrays Data and Desc in structure Package.Body;
***
### Default numbers:

```go
{
    PACK_TIME: 10,
    WAIT_TIME: 5,
    ROUTE_NUM: 3,
    BUFF_SIZE: 1 << 9, // 512B
    PID_SIZE: 1 << 3, // 8B
    SESSION_SIZE: 1 << 5, // 32B
    MAXSIZE_ADDRESS: 1 << 5, // 32B
    MAXSIZE_PACKAGE: 1 << 10 << 10, // 1MiB
    CLIENT_NAME_SIZE: 1 << 4, // 16B
}
```
* PACK_TIME (Save package by ID in 10 seconds. For broadcast);
* WAIT_TIME (Time for receive test connection package);
* ROUTE_NUM (The number of nodes through which the packet passes);
* BUFF_SIZE (Read package by NB part);
* PID_SIZE (Package ID in broadcast);
* SESSION_SIZE (Default: 32B. May be 16B if used "CRYPTO_SPEED");
* MAXSIZE_ADDRESS (Max address node);
* MAXSIZE_PACKAGE (The size of the package that can be read at a time);
* CLIENT_NAME_SIZE (Length client name in bytes);
***
### Default strings:
```go
{
    TEMPLATE: "0.0.0.0",
    CLIENT_NAME: "[CLIENT]",
    SEPARATOR: "\000\001\007\005\000",
    END_BYTES: "\000\005\007\001\000",
    NETWORK_NAME: "GENESIS",
    TITLE_CONNECT: "[TITLE:CONNECT]",
    TITLE_REDIRECT: "[TITLE:REDIRECT]",
    MODE_READ: "[MODE:READ]",
    MODE_SAVE: "[MODE:SAVE]",
    MODE_TEST: "[MODE:TEST]",
    MODE_REMV: "[MODE:REMV]",
    MODE_MERG: "[MODE:MERG]",
    MODE_DISTRIB: "[MODE:DISTRIB]",
    MODE_DECENTR: "[MODE:DECENTR]",
    MODE_READ_MERG: "[MODE:READ][MODE:MERG]",
    MODE_SAVE_MERG: "[MODE:SAVE][MODE:MERG]",
    MODE_DISTRIB_READ: "[MODE:DISTRIB][MODE:READ]",
    MODE_DISTRIB_SAVE: "[MODE:DISTRIB][MODE:SAVE]",
}
```
***
### Setting functions:
```go
func SettingsSet(settings SettingsType) []uint8 {}
func SettingsGet(key string) interface{} {}
```
***
### Network functions and methods:
```go
func NewNode(addr string) *Node {}
func (node *Node) Open() *Node {}
func (node *Node) Run(handleInit func(*Node), handleServer func(*Node, *Package), handleClient func(*Node, []string)) *Node {}
func (node *Node) ReadOnly(types ReadonlyType) *Node {}
func (node *Node) IsHidden(addr string) bool {}
func (node *Node) IsHandle(addr string) bool {}
func (node *Node) IsNode(addr string) bool {}
func (node *Node) IsMyAddress(addr string) bool {}
func (node *Node) IsMyHashname(hashname string) bool {}
func (node *Node) AppendToAccessList(access AccessType, addresses ...string) {}
func (node *Node) DeleteFromAccessList(addresses ...string) *Node {}
func (node *Node) InAccessList(access AccessType, addr string) bool {}
func (node *Node) InConnections(addr string) bool {}
func (node *Node) CreateRedirect(pack *Package) *Package {}
func (node *Node) Send(pack *Package) *Node {}
func (node *Node) SendRedirect(pack *Package) *Node {}
func (node *Node) SendInitRedirect(pack *Package) *Node {}
func (node *Node) SendToAll(pack *Package) *Node {}
func (node *Node) SendToAllWithout(pack *Package, sender string) *Node {}
func (node *Node) MergeConnect(addr string) *Node {}
func (node *Node) HiddenConnect(addr string) *Node {}
func (node *Node) ConnectToList(list ...interface{}) *Node {}
func (node *Node) Connect(data ...string) *Node {}
func (node *Node) Disconnect(addresses ...string) *Node {}
func (node *Node) GetConnections(relation RelationType) []string {}
func (node *Node) Close() *Node {}
```
***
### Cryptography functions and methods:
```go
func ParsePrivate(privData string) *rsa.PrivateKey {}
func ParsePublic(pubData string) *rsa.PublicKey {}
func Sign(priv *rsa.PrivateKey, data []byte) []byte {}
func Verify(pub *rsa.PublicKey, data, sign []byte) error {}
func HashSum(data []byte) []byte {}
func GenerateSessionKey(max int) []byte {}
func EncryptRSA(pub *rsa.PublicKey, data []byte) []byte {}
func DecryptRSA(priv *rsa.PrivateKey, data []byte) []byte {}
func EncryptAES(key, data []byte) []byte {}
func DecryptAES(key, data []byte) []byte {}
func StringPublic(pub *rsa.PublicKey) string {}
func StringPrivate(priv *rsa.PrivateKey) string {}
func (node *Node) StringPublic() string {}
func (node *Node) StringPrivate() string {}
func (node *Node) ParsePrivate(privData string) *Node {}
func (node *Node) SetPrivate(priv *rsa.PrivateKey) *Node {}
func (node *Node) GeneratePrivate(bits int) *Node {}
func (node *Node) DecryptRSA(data []byte) []byte {}
func (node *Node) Sign(data []byte) []byte {}
```
***
### Use in example F2F chat:
```go
package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/number571/gopeer"
)

const (
    TITLE_MESSAGE = "[TITLE:MESSAGE]"
    MODE_READ = "[MODE:READ]"
)

func init() {
    if len(os.Args) != 2 { panic("len args != 2") }
    gopeer.SettingsSet(gopeer.SettingsType{
        "IS_DECENTR": true,
        "HAS_CRYPTO": true,
        "HAS_ROUTING": true,
        "HAS_FRIENDS": true,
    })
}

func main() {
    node := gopeer.NewNode(os.Args[1]).GeneratePrivate(2048)
    node.Open().Run(handleInit, handleServer, handleClient).Close()
}

func handleInit(node *gopeer.Node) {
    node.ReadOnly(gopeer.ReadNode).ConnectToList(
        [2]string{":8080", "password"},
        [2]string{":7070", "password"},
        [2]string{":6060", "password"},
    )
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {
        case TITLE_MESSAGE:
            switch pack.Head.Mode {
                case MODE_READ:
                    message := strings.TrimLeft(pack.Body.Data[0], " ")
                    if message == "" { return }
                    fmt.Printf("[%s]: %s\n", pack.From.Address, message)
            }
    }
}

func handleClient(node *gopeer.Node, message []string) {
    switch message[0] {
        case "/exit": os.Exit(0)
        case "/whoami": fmt.Println("|", node.Hashname)  
        case "/hidden": node.HiddenConnect(strings.Join(message[1:], " "))
        case "/network": fmt.Println(node.GetConnections(gopeer.RelationAll))
        case "/send": 
            switch len(message[1:]) {
                case 0, 1: fmt.Println("[send] need > 0, 1 arguments")
                default: node.SendInitRedirect(&gopeer.Package{
                    To: gopeer.To{
                        Address: message[1],
                    },
                    Head: gopeer.Head{
                        Title: TITLE_MESSAGE,
                        Mode: MODE_READ,
                    },
                    Body: gopeer.Body{
                        Data: [gopeer.DATA_SIZE]string{strings.Join(message[2:], " ")},
                    },
                })
            }
        default: node.SendToAll(&gopeer.Package{
            Head: gopeer.Head{
                Title: TITLE_MESSAGE,
                Mode: MODE_READ,
            },
            Body: gopeer.Body{
                Data: [gopeer.DATA_SIZE]string{strings.Join(message, " ")},
            },
        })
    }
}
```
