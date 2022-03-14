package crypto

type Encrypter interface {
	Encrypt(msg []byte) []byte
}

type Decrypter interface {
	Decrypt(msg []byte) []byte
}

type Converter interface {
	String() string
	Bytes() []byte
	Type() string
	Size() uint64
}

type Hasher interface {
	Converter
}

type Cipher interface {
	Encrypter
	Decrypter
	Converter
}

type PubKey interface {
	Encrypter
	Converter
	Address() string
	Verify([]byte, []byte) bool
}

type PrivKey interface {
	Decrypter
	Converter
	Sign([]byte) []byte
	PubKey() PubKey
}

type Puzzle interface {
	Proof([]byte) uint64
	Verify([]byte, uint64) bool
}

type PRNG interface {
	String(uint64) string
	Bytes(uint64) []byte
	Uint64() uint64
}
