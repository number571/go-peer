package crypto

type iEncrypter interface {
	Encrypt(msg []byte) []byte
}

type iDecrypter interface {
	Decrypt(msg []byte) []byte
}

type iConverter interface {
	String() string
	Bytes() []byte
	Type() string
	Size() uint64
}

type IHasher interface {
	iConverter
}

type ICipher interface {
	iEncrypter
	iDecrypter
	iConverter
}

type IPubKey interface {
	iEncrypter
	iConverter
	Address() string
	Verify([]byte, []byte) bool
}

type IPrivKey interface {
	iDecrypter
	iConverter
	Sign([]byte) []byte
	PubKey() IPubKey
}

type IPuzzle interface {
	Proof([]byte) uint64
	Verify([]byte, uint64) bool
}

type IPRNG interface {
	String(uint64) string
	Bytes(uint64) []byte
	Uint64() uint64
}
