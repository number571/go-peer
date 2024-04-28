package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":{
			"method":"POST",
			"host":"hidden-echo-service",
			"path":"/echo",
			"body":"%s"
		}
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
	httpClient := http.Client{Timeout: (10 * time.Minute)}

	requestData := fmt.Sprintf(
		cRequestTemplate,
		pReceiver,
		base64.StdEncoding.EncodeToString(pMessage),
	)
	req, err := http.NewRequest(
		http.MethodPost,
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
