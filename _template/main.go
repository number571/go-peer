package main

import (
	"github.com/number571/gopeer"
)

const (
	ADDRESS = "ipv4:port"
	TITLE   = "TITLE"
)

func main() {
	listener := gopeer.NewListener(ADDRESS)
	listener.Open().Run(handleServer)
	defer listener.Close()

	// listener.NewClient(gopeer.GeneratePrivate(2048))
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
	client.HandleAction(TITLE, pack,
		func(client *gopeer.Client, pack *gopeer.Package) (set string) {

			return
		},
		func(client *gopeer.Client, pack *gopeer.Package) {

		},
	)
}
