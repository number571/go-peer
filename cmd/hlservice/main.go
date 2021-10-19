package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	lc "github.com/number571/gopeer/local"
	nt "github.com/number571/gopeer/network"
)

type Request struct {
	Host   string
	Path   string
	Method string
	Head   map[string]string
	Body   []byte
}

var (
	PrivKey  cr.PrivKey
	Services map[string]string
)

const (
	FileWithPrivKey  = "priv.key"
	FileWithPubKey   = "pub.key"
	FileWithServices = "services.json"
)

func init() {
	if !fileIsExist(FileWithPrivKey) {
		priv := cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))
		writeFile(FileWithPrivKey, []byte(priv.String()))
		writeFile(FileWithPubKey, []byte(priv.PubKey().String()))
	}
	spriv := string(readFile(FileWithPrivKey))
	PrivKey = cr.LoadPrivKeyByString(spriv)

	if !fileIsExist(FileWithServices) {
		services := make(map[string]string)
		services["route_service"] = "localhost:8080"
		writeFile(FileWithServices, serialize(services))
	}
	deserialize(readFile(FileWithServices), &Services)
	fmt.Println(Services)
}

func main() {
	fmt.Println("Service is listening...")
	client := lc.NewClient(PrivKey)
	nt.NewNode(client).
		Handle([]byte("/hls"), hlservice).Listen(":9571")
}

func hlservice(client *lc.Client, msg *lc.Message) []byte {
	request := &Request{}
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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()
	return data
}

func fileIsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return data
}

func writeFile(file string, data []byte) error {
	return ioutil.WriteFile(file, data, 0644)
}

func serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func deserialize(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
