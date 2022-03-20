package main

import (
	"fmt"
	"time"

	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/robfig/cron/v3"
)

func hmsDefaultInit() error {
	gConfig = config.NewConfig("hms.cfg")
	gDB = database.NewKeyValueDB("hms.db")

	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	scheduler := cron.New(cron.WithLocation(jakartaTime))
	scheduler.AddFunc(gConfig.CleanCron(), func() {
		gDB.Clean()
	})

	fmt.Printf("Service is listening [%s]...\n", gConfig.Address())
	return nil
}
