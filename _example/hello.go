package main

import (
	"encoding/json"
	"fmt"
	"github.com/number571/gopeer"
)

const (
	ADDRESS = ":8080"
	TITLE   = "TITLE"
)

var (
	anotherClient = new(gopeer.Client)
)

func main() {
	listener := gopeer.NewListener(ADDRESS)
	listener.Open().Run(handleServer)
	defer listener.Close()

	client := listener.NewClient(gopeer.GeneratePrivate(1024))
	anotherClient = listener.NewClient(gopeer.GeneratePrivate(1024))

	handleClient(client)
}

func handleClient(client *gopeer.Client) {
	dest := gopeer.NewDestination(&gopeer.Destination{
		Address: ADDRESS,
		Public:  anotherClient.Keys.Public,
	})

	client.Connect(dest)
	client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  TITLE,
			Option: gopeer.Get("OPTION_GET").(string),
		},
		Body: gopeer.Body{
			Data: "hello, world!",
		},
	})
	client.Disconnect(dest)
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
	client.HandleAction(TITLE, pack,
		func(client *gopeer.Client, pack *gopeer.Package) (set string) {
			fmt.Printf("[%s]: '%s'\n", pack.From.Sender.Hashname, pack.Body.Data)
			return set
		},
		func(client *gopeer.Client, pack *gopeer.Package) {
			// after receive result package
		},
	)
}

func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}
