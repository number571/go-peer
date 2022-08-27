// HLS - Hidden Lake Service
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/go-peer/modules/closer"
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
	fmt.Println()
	gLogger.Warning("Shutting down...")
	closer.CloseAll([]closer.ICloser{gServerHTTP, gNode})
}
