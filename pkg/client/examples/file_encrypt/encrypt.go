package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	payloadHead = 0x01
)

func encryptFile(client client.IClient, receiver asymmetric.IPubKey, filename string) {
	msgLimit := client.GetMessageLimit()
	headSize := hashing.CSHA256Size + (2 * encoding.CSizeUint64) + uint64(len(filename)) + 1
	if msgLimit <= headSize {
		panic("msgLimit <= headSize")
	}

	inputFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	hash := getFileHash(filename)
	count := getChunksCount(filename, msgLimit, headSize)

	buf := make([]byte, msgLimit-headSize)
	for i := 0; ; i++ {
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

		iBytes := encoding.Uint64ToBytes(uint64(i))
		countBytes := encoding.Uint64ToBytes(uint64(count))

		msg, err := client.EncryptPayload(
			receiver,
			payload.NewPayload64(
				payloadHead,
				bytes.Join(
					[][]byte{
						hash,
						iBytes[:],
						countBytes[:],
						[]byte(filename),
						{0x00},
						buf[:n],
					},
					[]byte{},
				),
			),
		)
		if err != nil {
			panic(err)
		}

		outputFile, err := os.Create(fmt.Sprintf("chunk_%d.enc", i))
		if err != nil {
			panic(err)
		}
		_, errW := outputFile.Write(msg.ToBytes())
		errC := outputFile.Close()
		if errW != nil || errC != nil {
			panic(errW)
		}
	}
}
