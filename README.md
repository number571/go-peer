# gopeer
> Framework for create decentralized networks. Version: 1.2.0s.

### Framework based applications:
* HiddenLake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "F2F network");

### Specifications:
1. Data transfer:
* Protocol: TCP;
* End to end encryption;
* Direct/Hidden connection;
2. Encryption:
* Symmetric algorithm: AES-CBC;
* Asymmetric algorithm: RSA-OAEP;
* Hash function: SHA256;

### Template:
```go
package main

import (
    gp "./gopeer"
)

func init() {
    gp.Set(gp.SettingsType{
        "NETW_NAME": "NET_TEMPLATE",
        "AKEY_SIZE": uint(3 << 10),
        "SKEY_SIZE": uint(1 << 5),
    })
}

func main() {
    node := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
    gp.NewListener(":8080", node).Run(handleFunc)
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
    WAIT_TIME uint
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
    WAIT_TIME: 5,       // seconds
    BUFF_SIZE: 4 << 10, // 4KiB
    PACK_SIZE: 2 << 20, // 2MiB
    MAPP_SIZE: 1024,    // elems
    AKEY_SIZE: 1024,    // bits
    SKEY_SIZE: 16,      // bytes
    RAND_SIZE: 16,      // bytes
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
gp.Set(gp.SettingsType{
    "AKEY_SIZE": uint(3 << 10),
    "SKEY_SIZE": uint(1 << 5),
})
```

### Network functions and methods:
```go
func NewListener(address string, client *Client) *Listener {}
func (listener *Listener) Run(handle func(*Client, *Package)) error {}
func NewClient(priv *rsa.PrivateKey) *Client {}
func (client *Client) Send(receiver *rsa.PublicKey, pack *Package) {}
func (client *Client) Request(receiver *rsa.PublicKey, pack *Package) error {}
func (client *Client) Response(pub *rsa.PublicKey) {}
func (client *Client) Connect(address string, handle func(*Client, *Package)) error {}
func (client *Client) Disconnect(address string) {}
func (client *Client) Public() *rsa.PublicKey {}
func (client *Client) Private() *rsa.PrivateKey {}
func (client *Client) StringPublic() string {}
func (client *Client) StringPrivate() string {}
func (client *Client) HashPublic() string {}
func (client *Client) AppendToF2F(pub *rsa.PublicKey) {}
func (client *Client) RemoveFromF2F(pub *rsa.PublicKey) {}
func (client *Client) InF2F(pub *rsa.PublicKey) bool {}
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
    mapping:     map[string]bool,
    publicKey:   *rsa.PublicKey,
    privateKey:  *rsa.PrivateKey,
    connections: map[net.Conn]string,
    actions:     map[string]chan bool,
    F2F:         {
        Enabled: bool,
        friends: map[string]*rsa.PublicKey,
    },
}
```
