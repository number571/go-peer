package main

import (
	"errors"
	"io"
	"os"
)

func getFileCount(filename string, msgLimit, headSize uint64) uint64 {
	inputFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	buf := make([]byte, msgLimit-headSize)

	result := uint64(0)
	for {
		n, err := inputFile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		if n == 0 {
			break
		}
		result += 1
	}

	return result
}
