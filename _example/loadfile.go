package main

import (
	"encoding/json"
	"fmt"
	"github.com/number571/gopeer"
)

var (
	ADDRESS1 = gopeer.Get("IS_CLIENT").(string)
	ADDRESS2 = ":8080"
)

var (
	anotherClient       = new(gopeer.Client)
	node2Key, node2Cert = gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
)

func main() {
	node1Key, node1Cert := gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
	listener1 := gopeer.NewListener(ADDRESS1)
	listener1.Open(&gopeer.Certificate{
		Cert: []byte(node1Cert),
		Key:  []byte(node1Key),
	}).Run(handleServer)
	defer listener1.Close()

	client := listener1.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

	listener2 := gopeer.NewListener(ADDRESS2)
	listener2.Open(&gopeer.Certificate{
		Cert: []byte(node2Cert),
		Key:  []byte(node2Key),
	}).Run(handleServer)
	defer listener2.Close()

	anotherClient = listener2.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

	anotherClient.Sharing.Perm = true
	anotherClient.Sharing.Path = "./"

	handleClient(client)
}

func handleClient(client *gopeer.Client) {
	dest := &gopeer.Destination{
		Address:     ADDRESS2,
		Certificate: []byte(node2Cert),
		Public:      anotherClient.Public(),
	}

	client.Connect(dest)
	client.LoadFile(dest, "archive.zip", "output.zip")
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {

}

func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}
