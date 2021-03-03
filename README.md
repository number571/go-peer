# gopeer
> Framework for create decentralized networks. Version: 1.2.6s.

### Framework based applications:
* Hidden Lake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "HL");
* Hidden Email Service: [github.com/number571/HES](https://github.com/number571/HES "HES");

### Research Article:
* The theory of the structure of hidden systems: [hiddensystems.pdf](https://github.com/Number571/gopeer/blob/master/hiddensystems.pdf "TSHS");

### Specifications:
* Type: Embedded;
* Protocol: TCP;
* Routing: Fill;
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
        "AKEY_SIZE": uint(1 << 10),
        "SKEY_SIZE": uint(1 << 4),
    })
}

func main() {
    gp.NewClient(
        gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)), 
        handleFunc,
    ).RunNode(":8080")
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
    SKEY_SIZE: 1 << 5,  // 2^5B = 32B
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
    "AKEY_SIZE": uint(1 << 10),
    "SKEY_SIZE": uint(1 << 4),
})
```

### Network functions and methods:
```go
func NewClient(priv *rsa.PrivateKey, handle func(*Client, *Package)) *Client {}
func NewPackage(title, data string) *Package {}
func (client *Client) Handle(title string, pack *Package, handle func(*Client, *Package) string) {}
func (client *Client) RunNode(address string) error {}
func (client *Client) Send(receiver *rsa.PublicKey, pack *Package, route []*rsa.PublicKey, ppsender *rsa.PrivateKey) (string, error) {}
func (client *Client) RoutePackage(receiver *rsa.PublicKey, pack *Package, route []*rsa.PublicKey, ppsender *rsa.PrivateKey) *Package {}
func (client *Client) Connect(address string) error {}
func (client *Client) Disconnect(address string) {}
func (client *Client) Encrypt(receiver *rsa.PublicKey, pack *Package) *Package {}
func (client *Client) Decrypt(pack *Package) *Package {}
func (client *Client) PublicKey() *rsa.PublicKey {}
func (client *Client) PrivateKey() *rsa.PrivateKey {}
func (f2f *friendToFriend) State() bool {}
func (f2f *friendToFriend) Switch() {}
func (f2f *friendToFriend) List() []rsa.PublicKey {}
func (f2f *friendToFriend) InList(pub *rsa.PublicKey) bool {}
func (f2f *friendToFriend) Append(pub *rsa.PublicKey) {}
func (f2f *friendToFriend) Remove(pub *rsa.PublicKey) {}
```

### Cryptography functions:
```go
func GenerateBytes(max uint) []byte {}
func GenerateKey(bits uint) *rsa.PrivateKey {}
func HashSum(data []byte) []byte {}
func HashPublicKey(pub *rsa.PublicKey) string {}
func BytesToPrivateKey(privData []byte) *rsa.PrivateKey {}
func BytesToPublicKey(pubData []byte) *rsa.PublicKey {}
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {}
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {}
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
func SerializePackage(pack *Package) string {}
func DeserializePackage(jsonData string) *Package {}
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

### Client structure:
```go
{
    handle:      func(*Client, *Package)
    mutex:       sync.Mutex,
    privateKey:  *rsa.PrivateKey,
    mapping:     map[string]bool,
    connections: map[net.Conn]string,
    actions:     map[string]chan bool,
    F2F:         {
        mutex:   sync.Mutex,
        enabled: bool,
        friends: map[string]*rsa.PublicKey,
    },
}
```
