package puzzle

type IPuzzle interface {
	ProofBytes([]byte) uint64
	VerifyBytes([]byte, uint64) bool
}
