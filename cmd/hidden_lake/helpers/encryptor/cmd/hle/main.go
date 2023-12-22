package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/app"
)

func main() {
	app, err := app.InitApp(".", "./priv.key")
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
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	<-shutdown
}
