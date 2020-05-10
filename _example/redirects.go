package main

/*
      A -> F
==================
   A--------->B
             /|
   /--------- |
   |          |
   v          v
   C<-------->D
   |          |
   v          v
   \--------->E
              |
              v
              F
==================
*/

import (
    "time"
    "encoding/json"
    "fmt"
    "github.com/number571/gopeer"
)

func init() {
    gopeer.Set(gopeer.SettingsType{
        "KEY_SIZE": uint64(1 << 10),
    })
}

var (
    ADDRESS1 = gopeer.Get("IS_CLIENT").(string)
    ADDRESS6 = gopeer.Get("IS_CLIENT").(string)
)

const (
    ADDRESS2 = ":6060"
    ADDRESS3 = ":7070"
    ADDRESS4 = ":8080"
    ADDRESS5 = ":9090"
    TITLE    = "TITLE"
)

var (
    node2Key, node2Cert = gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
    node3Key, node3Cert = gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
    node4Key, node4Cert = gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
    node5Key, node5Cert = gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
    node6Key, node6Cert = gopeer.GenerateCertificate(gopeer.Get("NETWORK").(string), gopeer.Get("KEY_SIZE").(uint16))
    client2             = new(gopeer.Client)
    client3             = new(gopeer.Client)
    client4             = new(gopeer.Client)
    client5             = new(gopeer.Client)
    client6             = new(gopeer.Client)
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
    client2 = listener2.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

    listener3 := gopeer.NewListener(ADDRESS3)
    listener3.Open(&gopeer.Certificate{
        Cert: []byte(node3Cert),
        Key:  []byte(node3Key),
    }).Run(handleServer)
    defer listener3.Close()
    client3 = listener3.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

    listener4 := gopeer.NewListener(ADDRESS4)
    listener4.Open(&gopeer.Certificate{
        Cert: []byte(node4Cert),
        Key:  []byte(node4Key),
    }).Run(handleServer)
    defer listener4.Close()
    client4 = listener4.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

    listener5 := gopeer.NewListener(ADDRESS5)
    listener5.Open(&gopeer.Certificate{
        Cert: []byte(node5Cert),
        Key:  []byte(node5Key),
    }).Run(handleServer)
    defer listener5.Close()
    client5 = listener5.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

    listener6 := gopeer.NewListener(ADDRESS6)
    listener6.Open(&gopeer.Certificate{
        Cert: []byte(node6Cert),
        Key:  []byte(node6Key),
    }).Run(handleServer)
    defer listener6.Close()
    client6 = listener6.NewClient(gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16)))

    // client6.Sharing.Perm = true
    // client6.Sharing.Path = "./"

    fmt.Println("A:", client.Hashname())
    fmt.Println("B:", client2.Hashname())
    fmt.Println("C:", client3.Hashname())
    fmt.Println("D:", client4.Hashname())
    fmt.Println("E:", client5.Hashname())
    fmt.Println("F:", client6.Hashname())
    fmt.Println()

    handleClient(client)
}

func handleClient(client *gopeer.Client) {
    dest := &gopeer.Destination{
        Address:     ADDRESS2,
        Certificate: []byte(node2Cert),
        Public:      client2.Public(),
    }

    client.Connect(dest)
    client3.Connect(dest)
    client4.Connect(dest)

    dest2 := &gopeer.Destination{
        Address:     ADDRESS4,
        Certificate: []byte(node4Cert),
        Public:      client4.Public(),
    }

    client3.Connect(dest2)

    dest3 := &gopeer.Destination{
        Address:     ADDRESS3,
        Certificate: []byte(node3Cert),
        Public:      client3.Public(),
    }

    dest4 := &gopeer.Destination{
        Address:     ADDRESS4,
        Certificate: []byte(node4Cert),
        Public:      client4.Public(),
    }

    client5.Connect(dest3)
    client5.Connect(dest4)

    dest5 := &gopeer.Destination{
        Address:     ADDRESS5,
        Certificate: []byte(node5Cert),
        Public:      client5.Public(),
    }
    client6.Connect(dest5)

    destFinal := &gopeer.Destination{
        Receiver: client6.Public(),
    }

    client.Connect(destFinal)
    hash := client6.Hashname()

    // client.LoadFile(destFinal, "archive.zip", "output.zip")

    for i := 0; i < 10; i++ {
        client.SendTo(destFinal, &gopeer.Package{
            Head: gopeer.Head{
                Title:  TITLE,
                Option: gopeer.Get("OPTION_GET").(string),
            },
            Body: gopeer.Body{
                Data: fmt.Sprintf("hello, world! [%d]", i),
            },
        })
        select {
        case <-client.Connections[hash].Action:
            // pass
        case <-time.After(time.Duration(5) * time.Second):
            fmt.Println("ERROR")
        }
    }
}

func handleServer(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(TITLE, pack,
        func(client *gopeer.Client, pack *gopeer.Package) (set string) {
            fmt.Printf("[%s->%s]: '%s'\n", pack.From.Sender.Hashname, client.Hashname(), pack.Body.Data)
            return set
        },
        func(client *gopeer.Client, pack *gopeer.Package) {
            client.Connections[pack.From.Sender.Hashname].Action <- true
        },
    )
}

func printJSON(data interface{}) {
    jsonData, _ := json.MarshalIndent(data, "", "\t")
    fmt.Println(string(jsonData))
}
