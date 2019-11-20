package main

import (
    "os"
    "fmt"
    "bufio"
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
        "IS_DISTRIB": true,
        "HAS_CRYPTO": true,
        "HAS_ROUTING": true,
    })
}

func main() {
    node := gopeer.NewNode(os.Args[1]).GeneratePrivate(2048)
    node.Open().Run(handleServer, handleClient).Close()
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

func handleClient(node *gopeer.Node) {
    node.ReadOnly(gopeer.ReadNode).ConnectToList(":8080",":7070",":6060")
    for {
        handleCLI(node, strings.Split(inputString(), " "))
    }
}

func inputString() string {
    msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    return strings.Replace(msg, "\n", "", -1)
}

func handleCLI(node *gopeer.Node, message []string) {
    switch message[0] {
        case "/exit": os.Exit(0)
        case "/whoami": fmt.Println("|", node.Hashname)  
        case "/connect": node.MergeConnect(strings.Join(message[1:], " "))
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
