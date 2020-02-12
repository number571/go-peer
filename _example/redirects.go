package main

// For version 1.1.0s;

/*
    A -> E
    A----------B
              /|
    /--------- |
    |          |
    C----------D
    |          |
    \----------E
*/

import (
    "encoding/json"
    "fmt"
    "github.com/number571/gopeer"
)

const (
    ADDRESS1 = ":5050"
    ADDRESS2 = ":6060"
    ADDRESS3 = ":7070"
    ADDRESS4 = ":8080"
    ADDRESS5 = ":9090"
    TITLE   = "TITLE"
)

var (
    node2Key, node2Cert = gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
    node3Key, node3Cert = gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
    node4Key, node4Cert = gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
    node5Key, node5Cert = gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), 1024)
    client2 = new(gopeer.Client)
    client3 = new(gopeer.Client)
    client4 = new(gopeer.Client)
    client5 = new(gopeer.Client)
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
    client2 = listener2.NewClient(gopeer.GeneratePrivate(1024))

    listener3 := gopeer.NewListener(ADDRESS3)
    listener3.Open(&gopeer.Certificate{
        Cert: []byte(node3Cert),
        Key:  []byte(node3Key),
    }).Run(handleServer)
    defer listener3.Close()
    client3 = listener3.NewClient(gopeer.GeneratePrivate(1024))

    listener4 := gopeer.NewListener(ADDRESS4)
    listener4.Open(&gopeer.Certificate{
        Cert: []byte(node4Cert),
        Key:  []byte(node4Key),
    }).Run(handleServer)
    defer listener4.Close()
    client4 = listener4.NewClient(gopeer.GeneratePrivate(1024))

    listener5 := gopeer.NewListener(ADDRESS5)
    listener5.Open(&gopeer.Certificate{
        Cert: []byte(node5Cert),
        Key:  []byte(node5Key),
    }).Run(handleServer)
    defer listener5.Close()
    client5 = listener5.NewClient(gopeer.GeneratePrivate(1024))

    client5.SetSharing(true, "./")

    fmt.Println(client.Hashname, client2.Hashname, client3.Hashname, client4.Hashname, client5.Hashname)
    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS2,
        Certificate: []byte(node2Cert),
        Public:  client2.Keys.Public,
    })

    client.Connect(dest)
    client3.Connect(dest)
    client4.Connect(dest)

    dest2 := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS4,
        Certificate: []byte(node4Cert),
        Public:  client4.Keys.Public,
    })

    client3.Connect(dest2)

    dest3 := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS5,
        Certificate: []byte(node5Cert),
        Public:  client5.Keys.Public,
    })

    client3.Connect(dest3)
    client4.Connect(dest3)

    destFinal := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS2,
        Certificate: []byte(node2Cert),
        Public:  client2.Keys.Public,
        Receiver: client5.Keys.Public,
    })

    client.Connect(destFinal)
    client.SendTo(destFinal, &gopeer.Package{
        Head: gopeer.Head{
            Title: TITLE,
            Option: gopeer.Get("OPTION_GET").(string),
        },
        Body: gopeer.Body{
            Data: "hello, world!",
        },
    })

    destFinal2 := gopeer.NewDestination(&gopeer.Destination{
        Address: ADDRESS4,
        Certificate: []byte(node4Cert),
        Public:  client4.Keys.Public,
        Receiver: client.Keys.Public,
    })
    client5.SendTo(destFinal2, &gopeer.Package{
        Head: gopeer.Head{
            Title: TITLE,
            Option: gopeer.Get("OPTION_GET").(string),
        },
        Body: gopeer.Body{
            Data: "hello, world!",
        },
    })
    fmt.Scanln()
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(TITLE, pack,
        func(client *gopeer.Client, pack *gopeer.Package) (set string) {
            fmt.Printf("[%s->%s]: '%s'\n", pack.From.Sender.Hashname, client.Hashname, pack.Body.Data)
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
