package base64

import (
	"errors"
	"math"
)

func GetSizeInBase64(pBytesNum uint64) (uint64, error) {
	if pBytesNum < 2 {
		return 0, errors.New("pBytesNum < 2")
	}
	// base64 encoding bytes with add 1/4 bytes of original
	// (-2) is a '=' characters in the suffix of encoding bytes
	return pBytesNum - uint64(math.Ceil(float64(pBytesNum)/4)) - 2, nil
}
