package main

import (
	"fmt"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	gp "github.com/number571/go-peer/settings"
)

func main() {
	client1 := lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))
	client2 := lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))

	msg := lc.NewMessage([]byte("header"), []byte("hello, world!"), 0)
	encmsg := client1.Encrypt(client2.PubKey(), msg)

	decmsg := client2.Decrypt(encmsg)

	fmt.Println(string(decmsg.Body.Data))
}
