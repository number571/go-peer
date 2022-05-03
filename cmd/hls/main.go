// HLS - Hidden Lake Service
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/settings"
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

	// set response route
	gNode.WithResponseRouter(func(node network.INode) []crypto.IPubKey {
		randSizeRoute := crypto.NewPRNG().Uint64() % cSizeRoute
		return local.NewSelector(nodesInOnline(node)).
			Shuffle().
			Return(randSizeRoute)
	})

	// turn on pseudo packages
	gNode.Pseudo().Switch(true)

	// turn on f2f mode
	gNode.F2F().Switch(gConfig.F2F().Status())
	for _, pubKey := range gConfig.F2F().PubKeys() {
		gNode.F2F().Append(pubKey)
	}

	// turn on online checker
	isOnline := gConfig.OnlineChecker().Status()
	gNode.Online().Switch(isOnline)
	gNode.Checker().Switch(isOnline)
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

	// network checker
	go func() {
		sett := gNode.Client().Settings()
		tchk := time.Duration(sett.Get(settings.TimePing))

		for {
			time.Sleep(time.Second * tchk)
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

func nodesInOnline(node network.INode) []crypto.IPubKey {
	inOnline := []crypto.IPubKey{}
	for _, info := range gNode.Checker().ListWithInfo() {
		if !info.Online() {
			continue
		}
		inOnline = append(inOnline, info.PubKey())
	}
	return inOnline
}
