package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":"%s"
	}`

	cJsonDataTemplate = `{
        "method":"POST",
        "host":"hidden-lake-messenger",
        "path":"/push",
        "head":{
			"%s": "%s"
        },
        "body":"%s"
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
	replacer := strings.NewReplacer("\n", "", "\t", "", "\r", "", " ", "", "\"", "\\\"")

	pseudonym := "Bob"
	requestData := replacer.Replace(
		fmt.Sprintf(
			cJsonDataTemplate,
			hlm_settings.CHeaderPseudonym,
			pseudonym,
			base64.StdEncoding.EncodeToString(pMessage),
		),
	)

	req, err := http.NewRequest(
		http.MethodPut,
		"http://localhost:7572/api/network/request",
		bytes.NewBufferString(fmt.Sprintf(cRequestTemplate, pReceiver, requestData)),
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
