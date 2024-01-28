package utils

import "math"

func GetChunksCount(pBytesNum, pChunkSize uint64) uint64 {
	return uint64(math.Ceil(float64(pBytesNum) / float64(pChunkSize)))
}
