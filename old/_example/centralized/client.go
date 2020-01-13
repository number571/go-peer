package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "github.com/number571/gopeer"
)

const (
    TITLE_NETWORK = "[TITLE:NETWORK]"
    TITLE_MESSAGE = "[TITLE:MESSAGE]"
    MODE_READ = "[MODE:READ]"
    MODE_SAVE = "[MODE:SAVE]"
    SERVER_ADDR = ":8080"
    SEPARATOR = "\000\001\003"
)

func init() {
    gopeer.SettingsSet(gopeer.SettingsType{
        "IS_DECENTR": true,
        "HAS_CRYPTO": true,
        "HAS_ROUTING": true,
    })
}

func main() {
    node := gopeer.NewNode(gopeer.SettingsGet("CLIENT_NAME").(string)).GeneratePrivate(2048)
    node.Run(handleServer, handleClient)
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {
        case TITLE_MESSAGE:
            switch pack.Head.Mode {
                case MODE_READ:
                    fmt.Printf("[%s]: %s\n", pack.From.Address, pack.Body.Data[0])
            }
        case TITLE_NETWORK:
            switch pack.Head.Mode {
                case MODE_SAVE:
                    list := strings.Split(pack.Body.Data[0], SEPARATOR)
                    for _, addr := range list {
                        fmt.Println("|", addr)
                    }
            }
    }
}

func handleClient(node *gopeer.Node) {
    node.Connect(SERVER_ADDR)
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
        case "/hidden": 
            if len(message[1:]) != 1 { 
                fmt.Println("len args != 1") 
                return
            }
            node.HiddenConnect(message[1])

        case "/whoami":
            fmt.Println("|", node.Hashname)

        case "/network":
            node.Send(&gopeer.Package{
                To: gopeer.To{
                    Address: SERVER_ADDR,
                },
                Head: gopeer.Head{
                    Title: TITLE_NETWORK,
                    Mode: MODE_READ,
                },
            })

        case "/send":
            if len(message[1:]) < 2 { 
                fmt.Println("len args < 2")
                return
            }
            node.SendInitRedirect(&gopeer.Package{
                To: gopeer.To{
                    Address: message[1],
                },
                Head: gopeer.Head{
                    Title: TITLE_MESSAGE,
                    Mode: MODE_READ,
                },
                Body: gopeer.Body{
                    Data: [gopeer.DATA_SIZE]string{
                        strings.Join(message[2:], " "),
                    },
                },
            })
    }
}
