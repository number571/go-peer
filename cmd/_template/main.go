package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	lc "github.com/number571/gopeer/local"
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
	client := lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))
	nt.NewNode(client).
		Handle([]byte("/msg"), msgRoute).Listen(":8080")
	// ...
}

func msgRoute(client *lc.Client, msg *lc.Message) []byte {
	hash := cr.LoadPubKey(msg.Head.Sender).Address()
	fmt.Printf("[%s] => '%s'\n", hash, msg.Body.Data)
	return msg.Body.Data
}
