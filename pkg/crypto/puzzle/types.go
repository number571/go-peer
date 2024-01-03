package puzzle

type IPuzzle interface {
	ProofBytes([]byte, uint64) uint64
	VerifyBytes([]byte, uint64) bool
}
