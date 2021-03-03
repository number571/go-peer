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
	gp.NewClient(
		gp.GenerateKey(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	).RunNode(":8080")
	// ...
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	// ...
}
