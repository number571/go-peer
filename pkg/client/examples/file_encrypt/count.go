package main

import (
	"math"
	"os"
)

func getChunksCount(filename string, msgLimit, headSize uint64) uint64 {
	stat, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	size := uint64(stat.Size())
	return uint64(math.Ceil(float64(size) / float64(msgLimit-headSize)))
}
