package main

import (
	gp "./gopeer"
)

func init() {
	gp.Set(gp.SettingsType{
		"AKEY_SIZE": uint(1 << 10),
		"SKEY_SIZE": uint(1 << 4),
	})
}

func main() {
	node := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	node.SetHandle(handleFunc)
	gp.NewNode(":8080", node).Run()
	// ...
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	// ...
}
