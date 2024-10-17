package puzzle

import (
	"bytes"
	"math"
	"math/big"
	"runtime"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cHashSizeInBits = hashing.CHasherSize * 8
)

var (
	_ IPuzzle = &sPoWPuzzle{}
)

type sPoWPuzzle struct {
	fDiff uint8
}

func NewPoWPuzzle(pDiff uint64) IPuzzle {
	if pDiff >= math.MaxUint8 {
		panic("diff >= 256")
	}
	return &sPoWPuzzle{
		fDiff: uint8(pDiff),
	}
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func (p *sPoWPuzzle) ProofBytes(pPackHash []byte, pParallel uint64) uint64 {
	var (
		closed = make(chan struct{})
		target = big.NewInt(1)
	)

	maxParallel := uint64(runtime.GOMAXPROCS(0))
	setParallel := pParallel
	if pParallel == 0 {
		setParallel = 1
	}
	if pParallel > maxParallel {
		setParallel = maxParallel
	}

	chNonce := make(chan uint64, setParallel)
	packHash := make([]byte, len(pPackHash))
	copy(packHash, pPackHash)

	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for i := uint64(0); i < setParallel; i++ {
		go func(i uint64) {
			intHash := big.NewInt(1)
			for nonce := i; nonce < math.MaxUint64; nonce += setParallel {
				select {
				case <-closed:
					return
				default:
					bNonce := encoding.Uint64ToBytes(nonce)
					hash := hashing.NewHasher(bytes.Join(
						[][]byte{packHash, bNonce[:]},
						[]byte{},
					)).ToBytes()
					intHash.SetBytes(hash)
					if intHash.Cmp(target) == -1 {
						chNonce <- nonce
						return
					}
				}
			}
		}(i)
	}

	result := <-chNonce
	close(closed)
	return result
}

// Verifies the work of the proof of work function.
func (p *sPoWPuzzle) VerifyBytes(pPackHash []byte, pNonce uint64) bool {
	var (
		intHash = big.NewInt(1)
		target  = big.NewInt(1)
	)

	packHash := make([]byte, len(pPackHash))
	copy(packHash, pPackHash)

	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	bNonce := encoding.Uint64ToBytes(pNonce)
	hash := hashing.NewHasher(bytes.Join(
		[][]byte{packHash, bNonce[:]},
		[]byte{},
	)).ToBytes()
	intHash.SetBytes(hash)
	return intHash.Cmp(target) == -1
}
