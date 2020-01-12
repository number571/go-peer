# gopeer
> Framework for create decentralized networks. Version: 1.0.0s.

### Specifications:
1. Data transfer:
* Direct
* End to end encryption
* Packages in blockchain
2. Encryption:
* Symmetric algorithm: AES256-CBC
* Asymmetric algorithm: RSA-OAEP
* Hash function: HMAC(SHA256)

### Framework based applications:
* Network [HiddenLake](https://github.com/number571/hidden-lake "hidden network");

### Template:
```go
package main

import (
    "github.com/number571/gopeer"
)

const (
    ADDRESS = "ipv4:port"
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
            
            return
        },
        func(client *gopeer.Client, pack *gopeer.Package) {
            
        },
    )
}
```

### Settings:
```go
type settingsStruct struct {
	TITLE_LASTHASH   string
	TITLE_CONNECT    string
	TITLE_DISCONNECT string
	OPTION_GET       string
	OPTION_SET       string
	NETWORK          string
	VERSION          string
	BUFFSIZE         uint16
	DIFFICULTY       uint8
	RETRY_NUMB       uint8
	RETRY_TIME       uint8
	TEMPLATE         string
	HMACKEY          string
	GENESIS          string
	NOISE            string
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
	TITLE_LASTHASH:   "[TITLE-LASTHASH]",
	TITLE_CONNECT:    "[TITLE-CONNECT]",
	TITLE_DISCONNECT: "[TITLE-DISCONNECT]",
	OPTION_GET:       "[OPTION-GET]", // Send
	OPTION_SET:       "[OPTION-SET]", // Receive
	NETWORK:          "NETWORK-NAME",
	VERSION:          "Version 1.0.0",
	BUFFSIZE:         512,
	DIFFICULTY:       15,
	RETRY_NUMB:       3,
	RETRY_TIME:       5, // Seconds
	TEMPLATE:         "0.0.0.0",
	HMACKEY:          "PASSWORD",
	GENESIS:          "[GENESIS-PACKAGE]",
	NOISE:            "1234567890ABCDEFGHIJKLMNOPQRSTUV",
}
```

### Network functions and methods:
```go
func NewListener(address string) *Listener {}
func (listener *Listener) NewClient(private *rsa.PrivateKey) *Client {}
func (listener *Listener) Open() *Listener {}
func (listener *Listener) Close() {}
func (listener *Listener) Run(handleServer func(*Client, *Package)) *Listener {}
func (client *Client) InConnections(hash string) bool {}
func (client *Client) IsConnected(hash string) bool {}
func (client *Client) HandleAction(title string, pack *Package, handleGet func(*Client, *Package) string, handleSet func(*Client, *Package)) bool {}
func (client *Client) Connect(dest *Destination) error {}
func (client *Client) Disconnect(dest *Destination) error {}
func (client *Client) Send(pack *Package) (*Package, error) {}
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
		Address: string,
	},
	To: {
		Receiver: {
			Hashname: string,
		},
		Address: string,
	},
	Body: {
		Data: string,
		Desc: {
			Rand: string,
			PrevHash: string,
			CurrHash: string,
			Sign: string,
			Nonce: uint64,
			Difficulty: uint8,
		},
	},
}
```

### Listener structure:
```go
{
	Address: {
		Ipv4: string,
		Port: string,
	},
	Setting: {
		Listen: net.Listen,
	},
	Clients: map[string]{
		Hashname: string,
		Keys: {
			Private: *rsa.PrivateKey,
			Public: *rsa.PublicKey,
		},
		Address: string,
		Connections: map[string]{
			Connected: bool,
			Session: []byte,
			PrevSession: []byte,
			Waiting: chan bool,
			Address: string,
			LastHash: string,
			Public: *rsa.PublicKey,
		},
	},
}
```
### Destination structure:
```go
{
	Address: string,
	Public: *rsa.Public,
}
```
