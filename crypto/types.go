package crypto

type IEncrypter interface {
	Encrypt(msg []byte) []byte
}

type IDecrypter interface {
	Decrypt(msg []byte) []byte
}

type IConverter interface {
	String() string
	Bytes() []byte
	Type() string
	Size() uint64
}
