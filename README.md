# gopeer
> Framework for create secure decentralized applications. Version: 1.3;

### Framework based applications
* Hidden Lake: [github.com/number571/HiddenLake](https://github.com/number571/HiddenLake "HL");
* Hidden Email Service: [github.com/number571/HES](https://github.com/number571/HES "HES");

### Research Article
* The theory of the structure of hidden systems: [hiddensystems.pdf](https://github.com/Number571/gopeer/blob/master/hiddensystems.pdf "TSHS");

### Specifications
* Type: Embedded;
* Protocol: TCP;
* Routing: Fill;
* Encryption: E2E;
* Symmetric algorithm: AES-CBC;
* Asymmetric algorithm: RSA-OAEP, RSA-PSS;
* Hash function: SHA256;

### Template
```go
package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	nt "github.com/number571/gopeer/network"
)

func init() {
	gp.Set(gp.SettingsType{
		"AKEY_SIZE": uint(1 << 10),
		"SKEY_SIZE": uint(1 << 4),
	})
}

func main() {
	fmt.Println("Node is listening...")
	nt.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))).
		Handle([]byte("/msg"), msgRoute).
		RunNode(":8080")
	// ...
}

func msgRoute(client *nt.Client, msg *nt.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
```

### Settings
```go
type SettingsType map[string]interface{}
type settingsStruct struct {
	END_BYTES []byte
	RET_BYTES []byte
	ROUTE_MSG []byte
	RETRY_NUM uint
	WAIT_TIME uint
	POWS_DIFF uint
	CONN_SIZE uint
	PACK_SIZE uint
	BUFF_SIZE uint
	MAPP_SIZE uint
	AKEY_SIZE uint
	SKEY_SIZE uint
	RAND_SIZE uint
}
```

### Default settings
```go
{
	END_BYTES: []byte("\000\005\007\001\001\007\005\000"),
	RET_BYTES: []byte("\000\001\007\005\005\007\001\000"),
	ROUTE_MSG: []byte("\000\001\002\003\004\005\006\007"),
	RETRY_NUM: 3,       // quantity
	WAIT_TIME: 20,      // seconds
	POWS_DIFF: 20,      // bits
	CONN_SIZE: 10,      // quantity
	PACK_SIZE: 8 << 20, // 8*(2^20)B = 8MiB
	BUFF_SIZE: 2 << 20, // 2*(2^20)B = 2MiB
	MAPP_SIZE: 2 << 10, // 2*(2^10)H = 88KiB
	AKEY_SIZE: 2 << 10, // 2*(2^10)b = 256B
	SKEY_SIZE: 1 << 5,  // 2^5B = 32B
	RAND_SIZE: 1 << 4,  // 2^4B = 16B
}
```

### Settings functions
```go
func Set(settings SettingsType) []uint8 {}
func Get(key string) interface{} {}
```

### Get/Set settings example
```go
var AKEY_SIZE = gopeer.Get("AKEY_SIZE").(uint)
gopeer.Set(gopeer.SettingsType{
	"AKEY_SIZE": uint(1 << 10),
	"SKEY_SIZE": uint(1 << 4),
})
```

### Network functions and methods
```go
func NewClient(priv crypto.PrivKey) *Client {}

func (client *Client) RunNode(address string) error {}
func (client *Client) Handle(title []byte, handle func(*Client, *Message) []byte) *Client {}

func (client *Client) Send(pack *Message, route *Route) ([]byte, error) {}
func (client *Client) RouteMessage(pack *Message, route *Route) *Message {}

func (client *Client) Connections() []string {}
func (client *Client) InConnections(address string) bool {}
func (client *Client) Connect(address ...string) []error {}
func (client *Client) Disconnect(address ...string) {}

func (client *Client) Encrypt(receiver crypto.PubKey, pack *Message) *Message {}
func (client *Client) Decrypt(pack *Message) *Message {}
func (client *Client) PubKey() crypto.PubKey {}
func (client *Client) PrivKey() crypto.PrivKey {}

func (f2f *friendToFriend) State() bool {}
func (f2f *friendToFriend) Switch() {}
func (f2f *friendToFriend) List() []*rsa.PublicKey {}
func (f2f *friendToFriend) InList(pub crypto.PubKey) bool {}
func (f2f *friendToFriend) Append(pubs ...crypto.PubKey) {}
func (f2f *friendToFriend) Remove(pubs ...crypto.PubKey) {}

func NewRoute(receiver crypto.PubKey) *Route {}
func (route *Route) WithSender(psender crypto.PrivKey) *Route {}
func (route *Route) WithRoutes(routes []crypto.PubKey) *Route {}

func NewMessage(title, data []byte) *Message {}
func (pack *Message) WithDiff(diff uint) *Message {}
```

### Cryptographic functions and methods
```go
func NewPrivKey(bits uint) PrivKey {}
func LoadPrivKey(pbytes []byte) PrivKey {}
func (key *PrivKeyRSA) PubKey() PubKey {}
func (key *PrivKeyRSA) Decrypt(msg []byte) []byte {}
func (key *PrivKeyRSA) Sign(msg []byte) []byte {}
func (key *PrivKeyRSA) Bytes() []byte {}
func (key *PrivKeyRSA) String() string {}
func (key *PrivKeyRSA) Type() string {}

func LoadPubKey(pbytes []byte) PubKey {}
func (key *PubKeyRSA) Encrypt(msg []byte) []byte {}
func (key *PubKeyRSA) Verify(msg []byte, sig []byte) bool {}
func (key *PubKeyRSA) Address() Address {}
func (key *PubKeyRSA) Bytes() []byte {}
func (key *PubKeyRSA) String() string {}
func (key *PubKeyRSA) Type() string {}

func NewCipher(key []byte) Cipher {}
func (cph *CipherAES) Encrypt(msg []byte) []byte {}
func (cph *CipherAES) Decrypt(msg []byte) []byte {}

func NewPuzzle(diff uint) Puzzle {}
func (puzzle *PuzzlePOW) Proof(packHash []byte) uint64 {}
func (puzzle *PuzzlePOW) Verify(packHash []byte, nonce uint64) bool {}

func RaiseEntropy(info, salt []byte, bits int) []byte {}

func HashSum(data []byte) []byte {}

func GenRand(max uint) []byte {}
```

### Encoding functions
```go
func Base64Encode(data []byte) string {}
func Base64Decode(data string) []byte {}
func ToBytes(num uint64) []byte {}
```

### Message structure
```go
{
	Head: {
		Title:   []byte,
		Rand:    []byte,
		Sender:  []byte,
		Session: []byte,
	},
	Body: {
		Data: []byte,
		Hash: []byte,
		Sign: []byte,
		Npow: uint64,
	},
}
```

### Client structure
```go
{
	mutex:       sync.Mutex,
	privateKey:  crypto.PrivKey,
	functions:   map[string]func(*Client, *Message) []byte,
	mapping:     map[string]bool,
	connections: map[string]net.Conn,
	actions:     map[string]chan []byte,
	F2F:         {
		mutex:   sync.Mutex,
		enabled: bool,
		friends: map[string]crypto.PubKey,
	},
}
```
