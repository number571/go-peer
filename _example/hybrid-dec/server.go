package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/number571/gopeer"
)

const (
    TITLE_MESSAGE = "[TITLE:MESSAGE]"
    MODE_READ = "[MODE:READ]"
)

func init() {
    if len(os.Args) != 2 { panic("len args != 2") }
    gopeer.SettingsSet(gopeer.SettingsType{
        "IS_DECENTR": true,
        "HAS_CRYPTO": true,
        "HAS_ROUTING": true,
    })
}

func main() {
    node := gopeer.NewNode(os.Args[1]).GeneratePrivate(2048)
    node.Open().Run(handleInit, handleServer, handleClient).Close()
}

func handleInit(node *gopeer.Node) {
    node.ConnectToList(
        map[string]string{
            ":8080": "password",
            ":7070": "password",
            ":6060": "password",
    })
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {
        case TITLE_MESSAGE:
            switch pack.Head.Mode {
                case MODE_READ:
                    message := strings.TrimLeft(pack.Body.Data[0], " ")
                    if message == "" { return }
                    fmt.Printf("[%s]: %s\n", pack.From.Address, message)
            }
    }
}

func handleClient(node *gopeer.Node, message []string) {
    switch message[0] {
        case "/exit": os.Exit(0)
        case "/whoami": fmt.Println("|", node.Hashname)  
        case "/hidden": node.HiddenConnect(strings.Join(message[1:], " "))
        case "/network": fmt.Println(node.GetConnections(gopeer.RelationAll))
        case "/send": 
            switch len(message[1:]) {
                case 0, 1: fmt.Println("[connect] need > 0, 1 arguments")
                default: node.SendInitRedirect(&gopeer.Package{
                    To: gopeer.To{
                        Address: message[1],
                    },
                    Head: gopeer.Head{
                        Title: TITLE_MESSAGE,
                        Mode: MODE_READ,
                    },
                    Body: gopeer.Body{
                        Data: [gopeer.DATA_SIZE]string{strings.Join(message[2:], " ")},
                    },
                })
            }
        default: node.SendToAll(&gopeer.Package{
            Head: gopeer.Head{
                Title: TITLE_MESSAGE,
                Mode: MODE_READ,
            },
            Body: gopeer.Body{
                Data: [gopeer.DATA_SIZE]string{strings.Join(message, " ")},
            },
        })
    }
}
