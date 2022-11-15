// HLS - Hidden Lake Service
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := initValues(); err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if err := gApp.Run(); err != nil {
		panic(err)
	}
	defer gApp.Close()

	fmt.Println("Service is running...")

	<-shutdown
	fmt.Println("\nShutting down...")
}
