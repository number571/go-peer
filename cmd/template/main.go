package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	nt "github.com/number571/gopeer/network"
)

func init() {
	gp.Set(gp.SettingsType{
		"AKEY_SIZE": uint(1 << 10),
		"SKEY_SIZE": uint(1 << 4),
	})
}

func main() {
	fmt.Println("Node is listening...")
	nt.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))).
		Handle([]byte("/msg"), msgRoute).
		RunNode(":8080")
	// ...
}

func msgRoute(client *nt.Client, msg *nt.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
