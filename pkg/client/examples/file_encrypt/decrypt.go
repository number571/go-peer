package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
)

func decryptFile(client client.IClient, decPrefix string) {
	encChunks := getChunks(".enc")
	filename := decryptChunks(client, encChunks)
	mergeDecryptedChunks(decPrefix, filename, len(encChunks))
}

func mergeDecryptedChunks(decPrefix, filename string, chunkCount int) {
	outputFile, err := os.Create(decPrefix + filename)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	for i := 0; i < chunkCount; i++ {
		chunkBytes, err := os.ReadFile(fmt.Sprintf("chunk_%d.dec", i))
		if err != nil {
			panic(err)
		}
		if _, err := outputFile.Write(chunkBytes); err != nil {
			panic(err)
		}
	}
}

func decryptChunks(client client.IClient, encChunks []string) string {
	headSize := hashing.CSHA256Size + (2 * encoding.CSizeUint64) + 1 + 1
	hashes := make([][]byte, 0, len(encChunks))

	filename := ""

	for _, ec := range encChunks {
		encBytes, err := os.ReadFile(ec)
		if err != nil {
			panic(err)
		}

		encMsg, err := message.LoadMessage(
			newSettings(client.GetPrivKey().GetSize()),
			encBytes,
		)
		if err != nil {
			panic(err)
		}

		_, pld, err := client.DecryptMessage(encMsg)
		if err != nil {
			panic(err)
		}

		if pld.GetHead() != payloadHead {
			panic("pld.GetHead() != payloadHead")
		}

		body := pld.GetBody()
		if len(body) < headSize {
			panic("len(body) < headSize")
		}

		hashes = append(hashes, body[:hashing.CSHA256Size])

		bufNum := [encoding.CSizeUint64]byte{}
		copy(bufNum[:], body[hashing.CSHA256Size:hashing.CSHA256Size+encoding.CSizeUint64])
		i := encoding.BytesToUint64(bufNum)

		copy(bufNum[:], body[hashing.CSHA256Size+encoding.CSizeUint64:hashing.CSHA256Size+2*encoding.CSizeUint64])
		count := encoding.BytesToUint64(bufNum)

		if count != uint64(len(encChunks)) {
			panic("count != uint64(len(encChunks))")
		}

		fileInfo := body[hashing.CSHA256Size+2*encoding.CSizeUint64:]
		index := bytes.Index(fileInfo, []byte{0x00})
		if index == -1 {
			panic("index == -1")
		}

		fn := string(fileInfo[:index])
		fb := fileInfo[index+1:]

		if filename == "" {
			filename = fn
		}
		if filename != fn {
			panic("filename != fn")
		}
		if hasNotWritableCharacters(fn) {
			panic("hasNotWritableCharacters(fn)")
		}

		if err := os.WriteFile(fmt.Sprintf("chunk_%d.dec", i), fb, 0644); err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(hashes); i++ {
		for j := i + 1; j < len(hashes)-1; j++ {
			if !bytes.Equal(hashes[i], hashes[j]) {
				panic("!bytes.Equal(hashes[i], hashes[j])")
			}
		}
	}

	return filename
}

func getChunks(suffix string) []string {
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	result := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "chunk_") && strings.HasSuffix(name, suffix) {
			result = append(result, name)
		}
	}
	return result
}

func hasNotWritableCharacters(pS string) bool {
	for _, c := range pS {
		if !unicode.IsGraphic(c) {
			return true
		}
	}
	return false
}
