package main

import (
	_ "embed"

	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	gopeer "github.com/number571/go-peer"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/app"
	"github.com/number571/go-peer/internal/flag"
)

func main() {
	args := os.Args[1:]

	if flag.GetBoolFlagValue(args, "version") {
		fmt.Println(gopeer.CVersion)
		return
	}

	app, err := app.InitApp(args, ".", "./priv.key")
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
			log.Fatal(err)
		}
	}()

	<-shutdown
}
