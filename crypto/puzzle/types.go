package puzzle

type IPuzzle interface {
	Proof([]byte) uint64
	Verify([]byte, uint64) bool
}
