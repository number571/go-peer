package main

import (
	"fmt"

	"github.com/number571/go-peer/cmd/hes/config"
	"github.com/number571/go-peer/cmd/hes/database"
)

func hesDefaultInit() error {
	gConfig = config.NewConfig("hes.cfg")
	gDB = database.NewKeyValueDB("hes.db")

	fmt.Printf("Service is listening [%s]...\n", gConfig.Address())
	return nil
}
