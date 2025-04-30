package main

import (
	"errors"
	"io"
	"os"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func decrypt(client client.IClient, outFilename, inFilename string) error {
	mapKeys := asymmetric.NewMapPubKeys(client.GetPrivKey().GetPubKey())

	infile, err := os.Open(inFilename) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() { _ = infile.Close() }()

	outfile, err := os.Create(outFilename) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() { _ = infile.Close() }()

	buf := make([]byte, client.GetMessageSize())
	for i := 0; ; i++ {
		n, err := infile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if uint64(n) != client.GetMessageSize() { //nolint:gosec
			return errors.New("uint64(n) != msgSize") //nolint:err113
		}
		_, chunk, err := client.DecryptMessage(mapKeys, buf[:n])
		if err != nil {
			return err
		}
		if _, err := outfile.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
