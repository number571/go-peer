package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/robfig/cron/v3"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

func hmsDefaultInit() error {
	var (
		initOnly bool
	)

	flag.BoolVar(&initOnly, "init", false, "run initialization only")
	flag.Parse()

	gSettings = hms_settings.NewSettings()
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

	if initOnly {
		os.Exit(0)
	}

	fmt.Printf("Service is listening [%s]...\n", gConfig.Address())
	return nil
}
