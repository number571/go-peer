package main

import (
    "strings"
    "github.com/number571/gopeer"
)

const (
    TITLE_NETWORK = "[TITLE:NETWORK]"
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
    node := gopeer.NewNode(SERVER_ADDR).GeneratePrivate(2048)
    node.Open().Run(handleInit, handleServer, handleClient).Close()
}

func handleInit(node *gopeer.Node) {
    node.ReadOnly(gopeer.ReadHandle)
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {
        case TITLE_NETWORK:
            switch pack.Head.Mode {
                case MODE_READ:
                    if !node.InConnections(pack.From.Address) { return }
                    var list []string 
                    for _, conn := range node.Network.Connections {
                        if conn.Hashname == pack.From.Hashname { continue }
                        list = append(list, conn.Hashname)
                    }
                    node.Send(&gopeer.Package{
                        To: gopeer.To{
                            Address: pack.From.Address,
                        },
                        Head: gopeer.Head{
                            Title: TITLE_NETWORK,
                            Mode: MODE_SAVE,
                        },
                        Body: gopeer.Body{
                            Data: [gopeer.DATA_SIZE]string{
                                strings.Join(list, SEPARATOR),
                            },
                        },
                    })
            }
    }
}

func handleClient(node *gopeer.Node, message []string) {

}
