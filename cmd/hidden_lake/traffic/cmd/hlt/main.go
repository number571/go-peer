package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/app"
	pkg_settings "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/settings"
	"github.com/number571/go-peer/internal/pprof"
)

func main() {
	pprof.RunPprofService(pkg_settings.CServiceName)

	app, err := app.InitApp(".")
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

	<-shutdown
}
