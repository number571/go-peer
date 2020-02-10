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
)

var (
    anotherClient = new(gopeer.Client)
)

func main() {
    listener1 := gopeer.NewListener(ADDRESS1)
    listener1.Open().Run(handleServer)
    defer listener1.Close()

    client := listener1.NewClient(gopeer.GeneratePrivate(1024))

    listener2 := gopeer.NewListener(ADDRESS2)
    listener2.Open().Run(handleServer)
    defer listener2.Close()

    anotherClient = listener2.NewClient(gopeer.GeneratePrivate(1024))

    anotherClient.Sharing.Perm = true
    anotherClient.Sharing.Path = "./"

    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS2,
        Public:  anotherClient.Keys.Public,
    })
    client.Connect(dest)
    client.LoadFile(dest, "archive.zip", "output.zip")
    client.Disconnect(dest)
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    
}

func printJSON(data interface{}) {
    jsonData, _ := json.MarshalIndent(data, "", "\t")
    fmt.Println(string(jsonData))
}
