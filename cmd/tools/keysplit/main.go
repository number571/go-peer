package main

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

func main() {
	if len(os.Args) != 2 {
		panic("len(os.Args) != 2")
	}

	keyBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var (
		privKey asymmetric.IPrivKey
		pubKey  asymmetric.IPubKey
	)

	privKey = asymmetric.LoadPrivKey(string(keyBytes))
	pubKey = asymmetric.LoadPubKey(string(keyBytes))

	if privKey != nil {
		kemPrivKey := encoding.HexEncode(privKey.GetKEMPrivKey().ToBytes())
		dsaPrivKey := encoding.HexEncode(privKey.GetDSAPrivKey().ToBytes())

		if err := os.WriteFile("priv.kem.key", []byte(kemPrivKey), 0600); err != nil {
			panic(err)
		}
		if err := os.WriteFile("priv.dsa.key", []byte(dsaPrivKey), 0600); err != nil {
			panic(err)
		}
		return
	}

	if pubKey != nil {
		kemPubKey := encoding.HexEncode(pubKey.GetKEMPubKey().ToBytes())
		dsaPubKey := encoding.HexEncode(pubKey.GetDSAPubKey().ToBytes())

		if err := os.WriteFile("pub.kem.key", []byte(kemPubKey), 0600); err != nil {
			panic(err)
		}
		if err := os.WriteFile("pub.dsa.key", []byte(dsaPubKey), 0600); err != nil {
			panic(err)
		}
		return
	}

	panic("got invalid key or file is not found")
}
