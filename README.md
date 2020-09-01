# gopeer
> Framework for create decentralized networks. Version: 1.2.2s.

### Framework based applications:
* HiddenLake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "F2F network");

### Specifications:
* Protocol: TCP;
* Encryption: E2E;
* Symmetric algorithm: AES-CBC;
* Asymmetric algorithm: RSA-OAEP, RSA-PSS;
* Hash function: SHA256;

### Template:
```go
package main

import (
    gp "./gopeer"
)

func init() {
    gp.Set(gp.SettingsType{
        "AKEY_SIZE": uint(3 << 10),
        "SKEY_SIZE": uint(1 << 5),
    })
}

func main() {
    node := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), handleFunc)
    gp.NewNode(":8080", node).Run()
    // ...
}

func handleFunc(client *gp.Client, pack *gp.Package) {
    // ...
}
```

### Settings:
```go
type SettingsType map[string]interface{}
type settingsStruct struct {
    END_BYTES string
    ROUTE_MSG string
    RETRY_NUM uint
    WAIT_TIME uint
    POWS_DIFF uint
    CONN_SIZE uint
    BUFF_SIZE uint
    PACK_SIZE uint
    MAPP_SIZE uint
    AKEY_SIZE uint
    SKEY_SIZE uint
    RAND_SIZE uint
}
```

### Default settings:
```go
{
    END_BYTES: "\000\005\007\001\001\007\005\000",
    ROUTE_MSG: "\000\001\002\003\004\005\006\007",
    RETRY_NUM: 3,       // quantity
    WAIT_TIME: 20,      // seconds
    POWS_DIFF: 20,      // bits
    CONN_SIZE: 10,      // quantity
    BUFF_SIZE: 2 << 20, // 2*(2^20)B = 2MiB
    PACK_SIZE: 4 << 20, // 4*(2^20)B = 4MiB
    MAPP_SIZE: 2 << 10, // 2*(2^10)H = 88KiB
    AKEY_SIZE: 2 << 10, // 2*(2^10)b = 256B
    SKEY_SIZE: 1 << 4,  // 2^4B = 16B
    RAND_SIZE: 1 << 4,  // 2^4B = 16B
}
```

### Settings functions:
```go
func Set(settings SettingsType) []uint8 {}
func Get(key string) interface{} {}
```

### Get/Set settings example:
```go
var AKEY_SIZE = gopeer.Get("AKEY_SIZE").(uint)
gopeer.Set(gopeer.SettingsType{
    "AKEY_SIZE": uint(3 << 10),
    "SKEY_SIZE": uint(1 << 5),
})
```

### Network functions and methods:
```go
// CREATE
func NewNode(address string, client *Client) *Listener {}
func NewClient(priv *rsa.PrivateKey, handle func(*Client, *Package)) *Client {}
func NewPackage(title, data string) *Package {}
// ACTIONS
func Handle(title string, client *Client, pack *Package, handle func(*Client, *Package) string) {}
func (listener *Listener) Run() error {}
func (listener *Listener) Close() error {}
func (client *Client) Send(receiver *rsa.PublicKey, pack *Package, route []*rsa.PublicKey, pseudoSender *Client) (string, error) {}
func (client *Client) Connect(address string) error {}
func (client *Client) Disconnect(address string) {}
func (client *Client) Encrypt(receiver *rsa.PublicKey, pack *Package) *Package {}
func (client *Client) Decrypt(pack *Package) *Package {}
// KEYS
func (client *Client) Public() *rsa.PublicKey {}
func (client *Client) Private() *rsa.PrivateKey {}
func (client *Client) StringPublic() string {}
func (client *Client) StringPrivate() string {}
func (client *Client) HashPublic() string {}
// F2F
func (client *Client) F2F() bool {}
func (client *Client) EnableF2F() {}
func (client *Client) DisableF2F() {}
func (client *Client) InF2F(pub *rsa.PublicKey) bool {}
func (client *Client) ListF2F() []rsa.PublicKey {}
func (client *Client) AppendF2F(pub *rsa.PublicKey) {}
func (client *Client) RemoveF2F(pub *rsa.PublicKey) {}
```

### Cryptography functions:
```go
func GenerateBytes(max uint) []byte {}
func GeneratePrivate(bits uint) *rsa.PrivateKey {}
func HashPublic(pub *rsa.PublicKey) string {}
func HashSum(data []byte) []byte {}
func ParsePrivate(privData string) *rsa.PrivateKey {}
func ParsePublic(pubData string) *rsa.PublicKey {}
func StringPrivate(priv *rsa.PrivateKey) string {}
func StringPublic(pub *rsa.PublicKey) string {}
func EncryptRSA(pub *rsa.PublicKey, data []byte) []byte {}
func DecryptRSA(priv *rsa.PrivateKey, data []byte) []byte {}
func Sign(priv *rsa.PrivateKey, data []byte) []byte {}
func Verify(pub *rsa.PublicKey, data, sign []byte) error {}
func EncryptAES(key, data []byte) []byte {}
func DecryptAES(key, data []byte) []byte {}
func ProofOfWork(packHash []byte, diff uint) uint64 {}
func ProofIsValid(packHash []byte, diff uint, nonce uint64) bool {}
```

### Additional functions:
```go
func EncodePackage(pack *Package) string {}
func DecodePackage(jsonData string) *Package {}
func Base64Encode(data []byte) string {}
func Base64Decode(data string) []byte {}
```

### Package structure:
```go
{
    Head: {
        Rand:    string,
        Title:   string,
        Sender:  string,
        Session: string,
    },
    Body: {
        Data: string,
        Hash: string,
        Sign: string,
        Npow: uint64,
    },
}
```

### Listener structure:
```go
{
    address: string,
    client:  *Client,
    listen:  net.Listener,
}
```

### Client structure:
```go
{
    mutex:       *sync.Mutex,
    privateKey:  *rsa.PrivateKey,
    mapping:     map[string]bool,
    connections: map[net.Conn]string,
    actions:     map[string]chan bool,
    f2f:         {
        enabled: bool,
        friends: map[string]*rsa.PublicKey,
    },
}
```
