package main

import (
	"errors"
	"io"
	"os"

	"github.com/number571/go-peer/pkg/client"
)

func encrypt(client client.IClient, outFilename, inFilename string) error {
	msgLimit := client.GetMessageLimit()

	infile, err := os.Open(inFilename)
	if err != nil {
		return err
	}
	defer infile.Close()

	outfile, err := os.Create(outFilename)
	if err != nil {
		return err
	}
	defer infile.Close()

	buf := make([]byte, msgLimit)
	for i := 0; ; i++ {
		n, err := infile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		kemPubKey := client.GetPrivKeyChain().GetKEncPrivKey().GetPubKey()
		chunk, err := client.EncryptMessage(kemPubKey, buf[:n])
		if err != nil {
			return err
		}
		if _, err := outfile.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
