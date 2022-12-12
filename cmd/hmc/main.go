// made by cryptohomochok
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/go-peer/pkg/closer"
	"github.com/number571/go-peer/pkg/types"
)

func main() {
	err := hmcDefaultInit()
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		gActions["load"].Do()
		for {
			cmd := newInputter("> ").String()
			f, ok := gActions[cmd]
			if !ok {
				fmt.Println("Undefined command")
				continue
			}
			f.Do()
		}
	}()

	<-shutdown
	fmt.Println("Shutting down...")
	closer.CloseAll([]types.ICloser{gDB})
}
