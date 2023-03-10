package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/number571/go-peer/internal/pprof"
)

func main() {
	pprof.RunPprofService(3, time.Second)

	app, err := initApp()
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if err := app.Run(); err != nil {
		panic(err)
	}
	defer func() {
		if err := app.Stop(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Service is running...")

	<-shutdown
	fmt.Println("\nShutting down...")
}
