package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":{
			"method":"GET",
			"host":"hidden-lake-filesharer",
			"path":"/list?page=0"
		}
	}`
)

func main() {
	receiver := "Bob"
	getListFiles(receiver)
}

func getListFiles(pReceiver string) {
	httpClient := http.Client{Timeout: time.Minute / 2}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8572/api/network/request",
		bytes.NewBufferString(fmt.Sprintf(cRequestTemplate, pReceiver)),
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
