package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	os.WriteFile("priv.key", pem.EncodeToMemory(privateKeyBlock), 0644)

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privKey.PublicKey),
	}
	os.WriteFile("pub.key", pem.EncodeToMemory(publicKeyBlock), 0644)
}
