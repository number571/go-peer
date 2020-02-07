# gopeer
> Framework for create decentralized networks. Version: 1.0.0s.

### Framework based applications:
* HiddenLake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "F2F network");

### Specifications:
1. Data transfer:
* Direct / Throw;
* File transfer supported;
* End to end encryption;
* Packages in blockchain;
2. Encryption:
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
type settingsStruct struct {
    TITLE_LASTHASH     string
    TITLE_CONNECT      string
    TITLE_DISCONNECT   string
    TITLE_FILETRANSFER string
    OPTION_GET         string
    OPTION_SET         string
    NETWORK            string
    VERSION            string
    PACKSIZE           uint32
    FILESIZE           uint32
    BUFFSIZE           uint16
    DIFFICULTY         uint8
    RETRY_NUMB         uint8
    RETRY_TIME         uint8
    TEMPLATE           string
    HMACKEY            string
    GENESIS            string
    NOISE              string
}
```

### Get/Set settings example:
```go
var OPTION_GET = gopeer.Get("OPTION_GET").(string)
gopeer.Set(gopeer.SettingsType{
    "NETWORK": "[HIDDEN-LAKE]",
    "VERSION": "[1.0.0s]",
    "HMACKEY": "9163571392708145",
    "GENESIS": "[GENESIS-LAKE]",
    "NOISE": "h19dlI#L9dkc8JA]1s-zSp,Nl/qs4;qf",
})
```

### Default settings:
```go
{
    TITLE_LASTHASH:     "[TITLE-LASTHASH]",
    TITLE_CONNECT:      "[TITLE-CONNECT]",
    TITLE_DISCONNECT:   "[TITLE-DISCONNECT]",
    TITLE_FILETRANSFER: "[TITLE-FILETRANSFER]",
    OPTION_GET:         "[OPTION-GET]", // Send
    OPTION_SET:         "[OPTION-SET]", // Receive
    NETWORK:            "NETWORK-NAME",
    VERSION:            "Version 1.0.0",
    PACKSIZE:           8 << 20, // 8MiB
    FILESIZE:           2 << 20, // 2MiB
    BUFFSIZE:           1 << 20, // 1MiB
    DIFFICULTY:         15,
    RETRY_NUMB:         2,
    RETRY_TIME:         5, // Seconds
    TEMPLATE:           "0.0.0.0",
    HMACKEY:            "PASSWORD",
    GENESIS:            "[GENESIS-PACKAGE]",
    NOISE:              "1234567890ABCDEFGHIJKLMNOPQRSTUV",
}
```

### Network functions and methods:
```go
func NewListener(address string) *Listener {}
func NewDestination(dest *Destination) *Destination
func (listener *Listener) NewClient(private *rsa.PrivateKey) *Client {}
func (listener *Listener) Open() *Listener {}
func (listener *Listener) Close() {}
func (listener *Listener) Run(handleServer func(*Client, *Package)) *Listener {}
func (client *Client) InConnections(hash string) bool {}
func (client *Client) HandleAction(title string, pack *Package, handleGet func(*Client, *Package) string, handleSet func(*Client, *Package)) bool {}
func (client *Client) LoadFile(dest *Destination, input string, output string) error
func (client *Client) Connect(dest *Destination) error {}
func (client *Client) Disconnect(dest *Destination) error {}
func (client *Client) SendTo(dest *Destination, pack *Package) (*Package, error) {}
```

### Cryptography functions:
```go
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
    Address: {
        Ipv4: string,
        Port: string,
    },
    Clients: map[string]{
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
            Session:   []byte,
            Address:   string,
            Public:    *rsa.PublicKey,
            PublicRecv *rsa.PublicKey,
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
