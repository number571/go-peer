package main

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
)

func fileHash(filename string) []byte {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err) //nolint:gocritic
	}

	return h.Sum(nil)
}
