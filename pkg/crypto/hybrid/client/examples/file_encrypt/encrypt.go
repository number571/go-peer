package main

import (
	"errors"
	"io"
	"os"

	"github.com/number571/go-peer/pkg/crypto/hybrid/client"
)

func encrypt(client client.IClient, outFilename, inFilename string) error {
	pldLimit := client.GetPayloadLimit()

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

	buf := make([]byte, pldLimit)
	for i := 0; ; i++ {
		n, err := infile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		pubKey := client.GetPrivKey().GetPubKey()
		chunk, err := client.EncryptMessage(pubKey, buf[:n])
		if err != nil {
			return err
		}
		if _, err := outfile.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
