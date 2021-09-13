package crypto

type Encrypter interface {
	Encrypt(msg []byte) []byte
}

type Decrypter interface {
	Decrypt(msg []byte) []byte
}

type Cipher interface {
	Encrypter
	Decrypter
}

type Address string
type PubKey interface {
	Encrypter
	Address() Address
	Bytes() []byte
	String() string
	Verify(msg []byte, sig []byte) bool
	Type() string
	Size() uint
}

type PrivKey interface {
	Decrypter
	Bytes() []byte
	String() string
	Sign(msg []byte) []byte
	Type() string
	PubKey() PubKey
}

type Puzzle interface {
	Proof(hash []byte) uint64
	Verify(hash []byte, nonce uint64) bool
}
