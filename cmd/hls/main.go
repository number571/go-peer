// HLS - Hidden Lake Service
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// read config, database, key
	err := hlsDefaultInit()
	if err != nil {
		gLogger.Error(err.Error())
		os.Exit(1)
	}

	// set handle functions
	gNode.Handle([]byte(cPatternHLS), routeHLS)

	// turn on f2f mode
	gNode.F2F().Switch(gConfig.F2F().Status())
	for _, pubKey := range gConfig.F2F().PubKeys() {
		gNode.F2F().Append(pubKey)
	}

	// turn on online checker
	gNode.Checker().Switch(gConfig.OnlineChecker().Status())
	for _, pubKey := range gConfig.OnlineChecker().PubKeys() {
		gNode.Checker().Append(pubKey)
	}

	// connect to open nodes
	for _, address := range gConfig.Connections() {
		if address == gConfig.Address().HLS() {
			gLogger.Warning(fmt.Sprintf("used own address '%s'", address))
			continue
		}
		err := gNode.Connect(address)
		if err != nil {
			gLogger.Warning(err.Error())
			continue
		}
		gLogger.Info(fmt.Sprintf("connected to '%s'", address))
	}

	go func() {
		timer := time.Second * 5
		for {
			time.Sleep(timer)
			for _, address := range gConfig.Connections() {
				if address == gConfig.Address().HLS() {
					continue
				}
				if gNode.InConnections(address) {
					continue
				}
				err := gNode.Connect(address)
				if err != nil {
					gLogger.Warning(err.Error())
					continue
				}
				gLogger.Info(fmt.Sprintf("connected to '%s'", address))
			}
		}
	}()

	// HTTP client
	http.HandleFunc("/", pageIndex)
	http.HandleFunc("/status", pageStatus)
	http.HandleFunc("/message", pageMessage)
	go func() {
		gLogger.Info(fmt.Sprintf("HTTP is listening [%s]...", gConfig.Address().HTTP()))
		err := http.ListenAndServe(gConfig.Address().HTTP(), nil)
		if err != nil {
			gLogger.Error(err.Error())
		}
	}()

	// if node in client mode
	// then run endless loop
	if gConfig.Address().HLS() == "" {
		gLogger.Info("HLS is listening...")
		select {}
	}

	// run node in server mode
	gLogger.Info(fmt.Sprintf("HLS is listening [%s]...", gConfig.Address().HLS()))
	err = gNode.Listen(gConfig.Address().HLS())
	if err != nil {
		gLogger.Error(err.Error())
		os.Exit(2)
	}
}
