# gopeer
> Framework for create decentralized networks. Version: 1.0.3s.

### Framework based applications:
* HiddenLake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "F2F network");

### Specifications:
1. Data transfer:
* Protocol: TCP;
* Direct / Throw;
* File transfer supported;
* End to end encryption;
2. Encryption:
* Protocol: TLS;
* Symmetric algorithm: AES256-CBC;
* Asymmetric algorithm: RSA-OAEP;
* Hash function: HMAC(SHA256);

### Template:
```go
package main

import (
    "github.com/number571/gopeer"
)

const (
    ADDRESS = ":8080"
    TITLE = "TITLE"
)

func main() {
    listener := gopeer.NewListener(ADDRESS)
    listener.Open().Run(handleServer)
    defer listener.Close()

    listener.NewClient(gopeer.GeneratePrivate(2048))
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(TITLE, pack, 
        func(client *gopeer.Client, pack *gopeer.Package) (set string) {
            return set
        },
        func(client *gopeer.Client, pack *gopeer.Package) {
            
        },
    )
}
```

### Settings:
```go
type SettingsType map[string]interface{}
type settingsStruct struct {
    TITLE_CONNECT      string
    TITLE_DISCONNECT   string
    TITLE_FILETRANSFER string
    OPTION_GET         string
    OPTION_SET         string
    SERVER_NAME        string
    IS_CLIENT          string
    END_BYTES          string
    TEMPLATE           string
    HMACKEY            string
    NETWORK            string
    VERSION            string
    PACK_SIZE          uint32
    FILE_SIZE          uint32
    BUFF_SIZE          uint32
    REMEMBER           uint16
    DIFFICULTY         uint8
    WAITING_TIME       uint8
    REDIRECT_QUAN      uint8
}
```

### Get/Set settings example:
```go
var OPTION_GET = gopeer.Get("OPTION_GET").(string)
gopeer.Set(gopeer.SettingsType{
    "NETWORK": "[HIDDEN-LAKE]",
    "VERSION": "[1.0.0s]",
    "HMACKEY": "9163571392708145",
})
```

### Default settings:
```go
{
    TITLE_CONNECT:      "[TITLE-CONNECT]",
    TITLE_DISCONNECT:   "[TITLE-DISCONNECT]",
    TITLE_FILETRANSFER: "[TITLE-FILETRANSFER]",
    OPTION_GET:         "[OPTION-GET]", // Send
    OPTION_SET:         "[OPTION-SET]", // Receive
    SERVER_NAME:        "GOPEER-FRAMEWORK",
    IS_CLIENT:          "[IS-CLIENT]", 
    END_BYTES:          "\000\000\000\005\007\001\000\000\000",
    TEMPLATE:           "0.0.0.0",
    HMACKEY:            "PASSWORD",
    NETWORK:            "NETWORK-NAME",
    VERSION:            "Version 1.0.0",
    PACK_SIZE:          8 << 20, // 8MiB
    FILE_SIZE:          2 << 20, // 2MiB
    BUFF_SIZE:          1 << 20, // 1MiB
    REMEMBER:           256, // hash packages
    DIFFICULTY:         15,
    WAITING_TIME:       5, // seconds
    REDIRECT_QUAN:      3,
}
```

### Settings functions:
```go
func Set(settings SettingsType) []uint8
func Get(key string) interface{}
```

### Network functions and methods:
```go
func NewListener(address string) *Listener {}
func NewDestination(dest *Destination) *Destination {}
func (listener *Listener) NewClient(private *rsa.PrivateKey) *Client {}
func (listener *Listener) Open(c *Certificate) *Listener {}
func (listener *Listener) Close() {}
func (listener *Listener) Run(handleServer func(*Client, *Package)) *Listener {}
func (client *Client) InConnections(hash string) bool {}
func (client *Client) HandleAction(title string, pack *Package, handleGet func(*Client, *Package) string, handleSet func(*Client, *Package)) bool {}
func (client *Client) LoadFile(dest *Destination, input string, output string) error
func (client *Client) Connect(dest *Destination) error {}
func (client *Client) Disconnect(dest *Destination) error {}
func (client *Client) SetSharing(perm bool, path string) {}
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {}
```

### Cryptography functions:
```go
func GenerateCertificate(name string, bits int) (string, string) {}
func GeneratePrivate(bits int) *rsa.PrivateKey {}
func ParsePrivate(privData string) *rsa.PrivateKey {}
func ParsePublic(pubData string) *rsa.PublicKey {}
func Sign(priv *rsa.PrivateKey, data []byte) []byte {}
func Verify(pub *rsa.PublicKey, data, sign []byte) error {}
func HashPublic(pub *rsa.PublicKey) string {}
func HashSum(data []byte) []byte {}
func HMAC(fHash func([]byte) []byte, data []byte, key []byte) []byte {}
func GenerateRandomIntegers(max int) []uint64 {}
func GenerateRandomBytes(max int) []byte {}
func EncryptRSA(pub *rsa.PublicKey, data []byte) []byte {}
func DecryptRSA(priv *rsa.PrivateKey, data []byte) []byte {}
func EncryptAES(key, data []byte) []byte {}
func DecryptAES(key, data []byte) []byte {}
func StringPublic(pub *rsa.PublicKey) string {}
func StringPrivate(priv *rsa.PrivateKey) string {}
func ProofOfWork(blockHash []byte, difficulty uint) uint64 {}
func NonceIsValid(blockHash []byte, difficulty uint, nonce uint64) bool {}
```

### Additional functions:
```go
func Base64Encode(data []byte) string {}
func Base64Decode(data string) []byte {}
func PackJSON(data interface{}) []byte {}
func UnpackJSON(jsonData []byte, data interface{}) interface{} {}
func ToBytes(num uint64) []byte {}
```

### Package structure:
```go
{
    Info: {
        Network: string,
        Version: string,
    },
    From: {
        Sender: {
            Hashname: string,
        },
        Hashname: string,
        Address:  string,
    },
    To: {
        Receiver: {
            Hashname: string,
        },
        Hashname: string,
        Address:  string,
    },
    Body: {
        Data: string,
        Desc: {
            Rand:       string,
            PrevHash:   string,
            CurrHash:   string,
            Sign:       string,
            Nonce:      uint64,
            Difficulty: uint8,
        },
    },
}
```

### Listener structure:
```go
{
    listen: net.Listen,
    handleFunc: func(*Client, *Package),
    Address: {
        Ipv4: string,
        Port: string,
    },
    Clients: map[string]{
        listener: *Listener,
        Hashname: string,
        Address:  string,
        Sharing: {
            Perm: bool,
            Path: string,
        },
        Keys: {
            Private: *rsa.PrivateKey,
            Public:  *rsa.PublicKey,
        },
        Connections: map[string]{
            connected:   bool,
            prevSession: []byte,
            waiting:     chan bool,
            lastHash:    string,
            transfer: {
                outputFile: string,
                isBlocked:  bool,
            },
            Session:    []byte,
            Address:    string,
            Relation:   net.Conn,
            Public:     *rsa.PublicKey,
            PublicRecv: *rsa.PublicKey,
        },
    },
}
```
### Destination structure:
```go
{
    Address:  string,
    Public:   *rsa.Public,
    Receiver: *rsa.Public,
}
```

### File Transfer structure:
```go
{
    Head: {
        Id:     uint32,
        Name:   string,
        IsNull: bool,
    },
    Body: {
        Hash: []byte,
        Data: []byte,
    },
}
```
