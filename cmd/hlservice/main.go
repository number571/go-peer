package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	lc "github.com/number571/gopeer/local"
	nt "github.com/number571/gopeer/network"
)

var (
	PrivKey  cr.PrivKey
	Services map[string]string
	Connects []string
)

const (
	FileWithPrivKey    = "priv.key"
	FileWithServices   = "services.json"
	FileWithConnects   = "connects.json"
	ServerAddressInRaw = "http://localhost:8080"
)

func init() {
	var initOnly bool

	flag.BoolVar(&initOnly, "init-only", false, "run initialization only")
	flag.Parse()

	if !fileIsExist(FileWithPrivKey) {
		priv := cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))
		writeFile(FileWithPrivKey, []byte(priv.String()))
		writeFile(FileWithPubKey, []byte(priv.PubKey().String()))
	}
	spriv := string(readFile(FileWithPrivKey))
	PrivKey = cr.LoadPrivKeyByString(spriv)

	if !fileIsExist(FileWithServices) {
		services := make(map[string]string)
		services[ServerAddressInHLS] = ServerAddressInRaw
		writeFile(FileWithServices, serialize(services))
	}
	deserialize(readFile(FileWithServices), &Services)

	if !fileIsExist(FileWithConnects) {
		connects := []string{"localhost:7070"}
		writeFile(FileWithConnects, serialize(connects))
	}
	deserialize(readFile(FileWithConnects), &Connects)

	if initOnly {
		os.Exit(0)
	}
}

func main() {
	fmt.Println("Service is listening...")
	client := lc.NewClient(PrivKey)
	node := nt.NewNode(client).Handle([]byte(HLS), hlservice)
	for _, conn := range Connects {
		err := node.Connect(conn)
		if err != nil {
			fmt.Println(err)
		}
	}
	node.Listen(":9571")
}

func hlservice(client *lc.Client, msg *lc.Message) []byte {
	request := new(Request)

	deserialize(msg.Body.Data, request)
	if request == nil {
		return nil
	}

	addr, ok := Services[request.Host]
	if !ok {
		return nil
	}

	req, err := http.NewRequest(
		request.Method,
		addr+request.Path,
		bytes.NewReader(request.Body),
	)
	if err != nil {
		return nil
	}

	for key, val := range request.Head {
		req.Header.Add(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return data
}
