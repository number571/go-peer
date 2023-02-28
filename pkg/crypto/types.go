package crypto

type IEncrypter interface {
	EncryptBytes(msg []byte) []byte
}

type IDecrypter interface {
	DecryptBytes(msg []byte) []byte
}
