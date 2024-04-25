package joiner

import (
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cAllocBytes = 1024
)

func NewBytesJoiner(pBytesSlice [][]byte) []byte {
	p := make([]byte, 0, cAllocBytes)
	for _, b := range pBytesSlice {
		p = append(
			p,
			payload.NewPayload(
				uint64(len(b)),
				b,
			).ToBytes()...,
		)
	}
	return p
}

func LoadBytesJoiner(pJoinerBytes []byte) ([][]byte, error) {
	p := make([][]byte, 0, cAllocBytes)

	for len(pJoinerBytes) != 0 {
		pld := payload.LoadPayload(pJoinerBytes)
		if pld == nil {
			return nil, ErrLoadPayload
		}

		pLen := pld.GetHead()
		pBytes := pld.GetBody()

		if pLen > uint64(len(pBytes)) {
			return nil, ErrInvalidLength
		}

		pbuf := make([]byte, len(pBytes[:pLen]))
		copy(pbuf, pBytes[:pLen])

		p = append(p, pbuf)
		pJoinerBytes = pBytes[pLen:]
	}

	return p, nil
}
