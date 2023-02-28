package entropy

type IEntropyBooster interface {
	BoostEntropy([]byte) []byte
}
