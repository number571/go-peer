package symmetric

type ICipher interface {
	EncryptBytes(pMsg []byte) []byte
	DecryptBytes(pMsg []byte) []byte
}
