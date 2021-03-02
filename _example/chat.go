package main

import (
	gp "./gopeer"
	"fmt"
	"strings"
	"bufio"
	"flag"
	"os"
	"crypto/rsa"
)

const (
	TITLE_MESSAGE = "TITLE_MESSAGE"
)

var (
	ADDRESS string
)

func init() {
	addrPtr := flag.String("open", "", "open node address")
	flag.Parse()
	ADDRESS = *addrPtr
}

func main() {
	var (
		message string
		splited []string
		receiver *rsa.PublicKey
		priv   = gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))
		client = gp.NewClient(priv, handleFunc)
	)
	go client.RunNode(ADDRESS)
	for {
		message = inputString("")
		splited = strings.Split(message, " ")
		switch splited[0] {
		case "/exit":
			os.Exit(0)
		case "/connect":
			if len(splited) != 2 {
				fmt.Println("error: len.splited != 2\n")
				continue
			}
			client.Connect(splited[1])
			fmt.Println("success: connect to node\n")
		case "/public":
			fmt.Printf("%s\n\n", gp.Base64Encode(gp.PublicKeyToBytes(client.PublicKey())))
		case "/receiver":
			if len(splited) != 2 {
				fmt.Println("error: len.splited != 2\n")
				continue
			}
			receiver = gp.BytesToPublicKey(gp.Base64Decode(splited[1]))
			fmt.Println("success: set receiver\n")
		case "/send":
			if len(splited) < 2 {
				fmt.Println("error: len.splited < 2\n")
				continue
			}
			if receiver == nil {
				fmt.Println("error: receiver is nil\n")
				continue
			}
			_, err := client.Send(
				receiver, 
				gp.NewPackage(TITLE_MESSAGE, strings.Join(splited[1:], " ")), 
				nil, 
				nil,
			)
			if err != nil {
				fmt.Println("error: send message\n")
				continue
			} 
			fmt.Println("success: message send\n")
		}
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	client.Handle(TITLE_MESSAGE, pack, getMessage)
}

func getMessage(client *gp.Client, pack *gp.Package) (set string) {
	publicBytes := gp.Base64Decode(pack.Head.Sender)
	hash := gp.Base64Encode(gp.HashSum(publicBytes))
	fmt.Printf("[%s] => '%s'\n", hash, pack.Body.Data)
	return "ok"
}

func inputString(before string) string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", 1)
}
