package main

import (
    "github.com/number571/gopeer"
)

func init() {
    gopeer.SettingsSet(gopeer.SettingsType{
        
    })
}

func main() {
    gopeer.NewNode("IPv4:Port").Open().Run(handleServer, handleClient).Close()
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {

    }
}

func handleClient(node *gopeer.Node) {
   
}
