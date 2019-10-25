package main

import (
    "github.com/number571/gopeer"
)

func init() {
    gopeer.SettingsSet(gopeer.SettingsType{
        
    })
}

func main() {
    gopeer.NewNode(":8080").Open().Run(handleInit, handleServer, handleClient).Close()
}

func handleInit(node *gopeer.Node) {
    
}

func handleServer(node *gopeer.Node, pack *gopeer.Package) {
    switch pack.Head.Title {

    }
}

func handleClient(node *gopeer.Node, message []string) {
    switch message[0] {
        
    }
}
