package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":{
			"method":"POST",
			"host":"hidden-lake-messenger",
			"path":"/push",
			"body":"%s"
		}
	}`
)

func main() {
	receiver := "Alice"
	message := "hello, world!"

	sendMessage(
		receiver,
		getRandomMessageType(message),
	)
}

func sendMessage(pReceiver string, pMessage []byte) {
	httpClient := http.Client{Timeout: time.Minute / 2}

	requestData := fmt.Sprintf(
		cRequestTemplate,
		pReceiver,
		base64.StdEncoding.EncodeToString(pMessage),
	)

	req, err := http.NewRequest(
		http.MethodPut,
		"http://localhost:7572/api/network/request",
		bytes.NewBufferString(requestData),
	)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}

func getRandomMessageType(pMessage string) []byte {
	if random.NewStdPRNG().GetBool() { // isText
		return bytes.Join(
			[][]byte{
				{hlm_settings.CIsText},
				[]byte(pMessage),
			},
			[]byte{},
		)
	}
	// isFile
	return bytes.Join(
		[][]byte{
			{hlm_settings.CIsFile},
			[]byte("example.txt"),
			{hlm_settings.CIsFile},
			[]byte(pMessage),
		},
		[]byte{},
	)
}
