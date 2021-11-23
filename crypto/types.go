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
	Size() uint
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
	Verify(msg []byte, sig []byte) bool
}

type PrivKey interface {
	Decrypter
	Converter
	Sign(msg []byte) []byte
	PubKey() PubKey
}

type Puzzle interface {
	Proof(hash []byte) uint64
	Verify(hash []byte, nonce uint64) bool
}
