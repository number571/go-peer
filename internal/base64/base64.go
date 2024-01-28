package base64

import "math"

func GetSizeInBase64(pBytesNum uint64) uint64 {
	return pBytesNum - uint64(math.Ceil(float64(pBytesNum)/4)) - 2
}
