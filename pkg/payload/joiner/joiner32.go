package joiner

import (
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cAllocBytes = 4096
)

func NewBytesJoiner32(pBytesSlice [][]byte) []byte {
	p := make([]byte, 0, cAllocBytes)
	for _, b := range pBytesSlice {
		p = append(
			p,
			payload.NewPayload32(
				uint32(len(b)), //nolint:gosec
				b,
			).ToBytes()...,
		)
	}
	return p
}

func LoadBytesJoiner32(pJoinerBytes []byte) ([][]byte, error) {
	p := make([][]byte, 0, cAllocBytes)

	for len(pJoinerBytes) != 0 {
		pld := payload.LoadPayload32(pJoinerBytes)
		if pld == nil {
			return nil, ErrLoadPayload
		}

		pLen := pld.GetHead()
		pBytes := pld.GetBody()

		if pLen > uint32(len(pBytes)) { //nolint:gosec
			return nil, ErrInvalidLength
		}

		pbuf := make([]byte, len(pBytes[:pLen]))
		copy(pbuf, pBytes[:pLen])

		p = append(p, pbuf)
		pJoinerBytes = pBytes[pLen:]
	}

	return p, nil
}
