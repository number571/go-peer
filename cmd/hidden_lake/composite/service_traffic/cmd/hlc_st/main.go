package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := initApp(".", "./priv.key")
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	closed := make(chan struct{})
	defer func() {
		cancel()
		<-closed
	}()

	go func() {
		defer func() { closed <- struct{}{} }()
		if err := app.Run(ctx); err != nil {
			panic(err)
		}
	}()

	<-shutdown
}
