package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/number571/go-peer/internal/api"
)

func main() {
	if len(os.Args) != 4 {
		panic("./producer [incoming-addr] [service-addr] [logger]")
	}

	incomingAddr := os.Args[1]
	serviceAddr := os.Args[2]
	logger := os.Args[3]

	http.HandleFunc("/traffic", trafficPage(serviceAddr, logger == "true"))
	_ = http.ListenAndServe(incomingAddr, nil)
}

func trafficPage(serviceAddr string, hasLog bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			printLog(hasLog, errors.New("r.Method != http.MethodPost"))
			api.Response(w, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		// get message from HLT
		msgStringAsBytes, err := io.ReadAll(r.Body)
		if err != nil {
			printLog(hasLog, err)
			api.Response(w, http.StatusConflict, "failed: read body")
			return
		}

		ret, res := pushMessageToService(serviceAddr, msgStringAsBytes, hasLog)
		api.Response(w, ret, res)
	}
}

func pushMessageToService(serviceAddr string, msgStringAsBytes []byte, hasLog bool) (int, string) {
	// build request to service
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		fmt.Sprintf("http://%s/push", serviceAddr),
		bytes.NewBuffer(msgStringAsBytes),
	)
	if err != nil {
		printLog(hasLog, err)
		return http.StatusNotImplemented, "failed: build request"
	}

	// send request to service
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		printLog(hasLog, err)
		return http.StatusBadRequest, "failed: bad request"
	}
	defer resp.Body.Close()

	// read response from service
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		printLog(hasLog, err)
		return http.StatusBadGateway, "failed: read body from service"
	}

	// read body of response
	if len(res) == 0 || res[0] == '!' {
		printLog(hasLog, errors.New("len(res) == 0 || res[0] == '!'"))
		return http.StatusForbidden, "failed: incorrect response from service"
	}

	return http.StatusOK, "success: push to service"
}

func printLog(hasLog bool, msg error) {
	if !hasLog {
		return
	}
	fmt.Println(msg)
}
