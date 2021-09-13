package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
)

func init() {
	gp.Set(gp.SettingsType{
		"AKEY_SIZE": uint(1 << 10),
		"SKEY_SIZE": uint(1 << 4),
	})
}

func main() {
	fmt.Println("Node is listening...")
	gp.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))).
		Handle([]byte("/msg"), msgRoute).
		RunNode(":8080")
	// ...
}

func msgRoute(client *gp.Client, pack *gp.Package) []byte {
	hash := cr.LoadPubKey(pack.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return pack.Body.Data
}
