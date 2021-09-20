package main

import (
	"fmt"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	lc "github.com/number571/gopeer/local"
)

func main() {
	client1 := lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))
	client2 := lc.NewClient(cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint)))

	msg := lc.NewMessage([]byte("header"), []byte("hello, world!"))
	encmsg := client1.Encrypt(client2.PubKey(), msg.WithDiff(0))

	decmsg := client2.Decrypt(encmsg)

	fmt.Println(string(decmsg.Body.Data))
}
