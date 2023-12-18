package crypto

type IEncrypter interface {
	EncryptBytes(pMsg []byte) []byte
}

type IDecrypter interface {
	DecryptBytes(pMsg []byte) []byte
}
