package main

import (
	gp "./gopeer"
	"fmt"
	"time"
	// "encoding/json"
)

const (
	TITLE_MESSAGE = "TITLE_MESSAGE"
	NODE_ADDRESS  = ":8080"
)

func main() {
	client := gp.NewClient(
		gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	node := gp.NewClient(
		gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)),
		handleFunc,
	)
	go node.RunNode(NODE_ADDRESS)
	time.Sleep(500 * time.Millisecond)

	client.Connect(NODE_ADDRESS)

	res, err := client.Send(
		node.Public(), 
		gp.NewPackage(TITLE_MESSAGE, "hello, world!"), 
		nil, 
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	client.Handle(TITLE_MESSAGE, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	// printJSON(pack)
	public := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(public), pack.Body.Data)
	return "ok"
}

// func printJSON(data interface{}) {
// 	jsonData, _ := json.MarshalIndent(data, "", "\t")
// 	fmt.Println(string(jsonData))
// }
