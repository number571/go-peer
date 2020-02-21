package main

import (
	"github.com/number571/gopeer"
)

const (
	ADDRESS = "ipv4:port"
	TITLE   = "TITLE"
)

func main() {
	key, cert := gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
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
			return
		},
		func(client *gopeer.Client, pack *gopeer.Package) {
		},
	)
	// ...
}
