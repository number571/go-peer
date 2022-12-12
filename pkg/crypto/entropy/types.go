package entropy

type IEntropy interface {
	Raise([]byte, []byte) []byte
}
