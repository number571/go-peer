package main

import (
	gp "./gopeer"
)

func init() {
	gp.Set(gp.SettingsType{
		"NETW_NAME": "NET_TEMPLATE",
		"AKEY_SIZE": uint(3 << 10),
		"SKEY_SIZE": uint(1 << 5),
	})
}

func main() {
	node := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	gp.NewListener(":8080", node).Run(handleFunc)
	// ...
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	// ...
}
