package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/cmd/hes/config"
	"github.com/number571/go-peer/cmd/hes/database"
	"github.com/robfig/cron/v3"
)

func hesDefaultInit() error {
	gConfig = config.NewConfig("hes.cfg")
	gDB = database.NewKeyValueDB("hes.db")

	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	scheduler := cron.New(cron.WithLocation(jakartaTime))
	scheduler.AddFunc("0 0 * * *", func() {
		gDB.Clean()
	})

	fmt.Printf("Service is listening [%s]...\n", gConfig.Address())
	return nil
}
