// HLS - Hidden Lake Service
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/go-peer/utils"
)

func main() {
	err := hlsDefaultInit()
	if err != nil {
		gLogger.Error(err.Error())
		os.Exit(1)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	fmt.Println("Shutting down...")
	utils.CloseAll([]utils.ICloser{gServerHTTP, gNode})
}
