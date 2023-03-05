package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/adapters/common"
	"github.com/number571/go-peer/internal/api"
)

func main() {
	if len(os.Args) != 3 {
		panic("./sender [port-incoming] [port-service]")
	}

	portIncoming, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	portService, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/traffic", trafficPage(portService))
	http.ListenAndServe(fmt.Sprintf(":%d", portIncoming), nil)
}

func trafficPage(portService int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			api.Response(w, 2, "failed: incorrect method")
			return
		}

		res, err := io.ReadAll(r.Body)
		if err != nil {
			api.Response(w, 3, "failed: read body")
			return
		}

		// convert message to service pattern
		_, err = api.Request(
			http.MethodPost,
			fmt.Sprintf("%s:%d/push", common.HostService, portService),
			res,
		)
		if err != nil {
			api.Response(w, 4, "failed: bad response")
			return
		}

		api.Response(w, 1, "success: push to service")
	}
}
