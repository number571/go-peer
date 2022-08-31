package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/go-peer/modules/closer"
	"github.com/number571/go-peer/modules/inputter"
)

func main() {
	err := hlmDefaultInit()
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// read-only mode
		if gChannelPubKey != nil {
			return
		}
		// another mode
		for {
			cmd := inputter.NewInputter("> ").String()
			f, ok := gActions[cmd]
			if !ok {
				sendActionDefault(cmd)
				continue
			}
			f.Do()
		}
	}()

	<-shutdown
	fmt.Println()
	gLogger.Warning("Shutting down...")
	closer.CloseAll([]closer.ICloser{gServerHTTP})
}
