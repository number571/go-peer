package main

import (
	"github.com/number571/gopeer"
)

const (
	ADDRESS = "ipv4:port"
	TITLE   = "TITLE"
)

func init() {
	gopeer.Set(gopeer.SettingsType{
		"NETWORK": "GOPEER-NETWORK",
		"VERSION": "template 1.0.0",
		"KEY_SIZE": uint64(1 << 10),
	})
}

func main() {
	key, cert := gopeer.GenerateCertificate(
		gopeer.Get("NETWORK").(string), 
		gopeer.Get("KEY_SIZE").(uint16),
	)
	listener := gopeer.NewListener(ADDRESS)
	listener.Open(&gopeer.Certificate{
		Cert: []byte(cert),
		Key:  []byte(key),
	}).Run(handleServer)
	defer listener.Close()
	// ...
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
	client.HandleAction(TITLE, pack,
		func(client *gopeer.Client, pack *gopeer.Package) (set string) {
			return set
		},
		func(client *gopeer.Client, pack *gopeer.Package) {
		},
	)
	// ...
}
