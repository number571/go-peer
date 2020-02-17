package main

import (
    "encoding/json"
    "fmt"
    "github.com/number571/gopeer"
)

var (
    ADDRESS1 = gopeer.Get("IS_CLIENT").(string)
)

const (
    ADDRESS2 = ":8080"
    TITLE    = "TITLE"
)

var (
    anotherClient       = new(gopeer.Client)
    node2Key, node2Cert = gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
)

func main() {
    node1Key, node1Cert := gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
    listener1 := gopeer.NewListener(ADDRESS1)
    listener1.Open(&gopeer.Certificate{
        Cert: []byte(node1Cert),
        Key:  []byte(node1Key),
    }).Run(handleServer)
    defer listener1.Close()

    client := listener1.NewClient(gopeer.GeneratePrivate(1024))

    listener2 := gopeer.NewListener(ADDRESS2)
    listener2.Open(&gopeer.Certificate{
        Cert: []byte(node2Cert),
        Key:  []byte(node2Key),
    }).Run(handleServer)
    defer listener2.Close()

    anotherClient = listener2.NewClient(gopeer.GeneratePrivate(1024))

    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := &gopeer.Destination{
        Address:     ADDRESS2,
        Public:      anotherClient.Keys.Public,
    }

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
    fmt.Scanln()
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
