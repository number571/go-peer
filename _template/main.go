package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
)

func init() {
	gp.Set(gp.SettingsType{
		"AKEY_SIZE": uint(1 << 10),
		"SKEY_SIZE": uint(1 << 4),
	})
}

func main() {
	gp.NewClient(gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))).
		Handle("/msg", msgRoute).
		RunNode(":8080")
	// ...
}

func msgRoute(client *gp.Client, pack *gp.Package) []byte {
	hash := gp.HashPublicKey(gp.BytesToPublicKey(pack.Head.Sender))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
