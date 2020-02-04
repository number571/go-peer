package main

import (
    "fmt"
    "encoding/json"
    "github.com/number571/gopeer"
)

const (
    ADDRESS = ":8080"
    TITLE_MESSAGE = "MESSAGE"
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

    anotherClient.Sharing.Perm  = true
    anotherClient.Sharing.Path = "./"

    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := &gopeer.Destination{
        Address: ADDRESS,
        Public: anotherClient.Keys.Public,
    }
    client.Connect(dest)

    client.SendTo(dest, &gopeer.Package{
        Head: gopeer.Head{
            Title: TITLE_MESSAGE,
            Option: gopeer.Get("OPTION_GET").(string),
        },
        Body: gopeer.Body{
            Data: "begin",
        },
    })

    client.LoadFile(dest, "archive.zip", "output.zip")

    client.SendTo(dest, &gopeer.Package{
        Head: gopeer.Head{
            Title: TITLE_MESSAGE,
            Option: gopeer.Get("OPTION_GET").(string),
        },
        Body: gopeer.Body{
            Data: "end",
        },
    })

    client.Disconnect(dest)
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(TITLE_MESSAGE, pack, 
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
