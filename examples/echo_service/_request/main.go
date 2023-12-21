package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/random"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":"%s"
	}`

	cJsonDataTemplate = `{
        "method":"POST",
        "host":"hidden-echo-service",
        "path":"/echo",
        "head":{
			"%s": "%s",
            "Accept": "application/json"
        },
        "body":"%s"
	}`
)

func main() {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		diff := t2.Sub(t1)
		fmt.Println("Request took", diff)
	}()

	receiver := "Bob"
	message := "hello, world!"

	sendMessage(receiver, []byte(message))
}

func sendMessage(pReceiver string, pMessage []byte) {
	httpClient := http.Client{Timeout: time.Minute / 2}
	replacer := strings.NewReplacer("\n", "", "\t", "", "\r", "", " ", "", "\"", "\\\"")

	requestData := replacer.Replace(
		fmt.Sprintf(
			cJsonDataTemplate,
			hls_settings.CHeaderRequestId,
			random.NewStdPRNG().GetString(hls_settings.CHandleRequestIDSize),
			base64.StdEncoding.EncodeToString(pMessage),
		),
	)
	req, err := http.NewRequest(
		http.MethodPost,
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
