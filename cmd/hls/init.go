package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/utils"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"
	"github.com/robfig/cron/v3"
)

func hlsDefaultInit() error {
	var (
		initOnly bool
	)

	flag.BoolVar(&initOnly, "init-only", false, "run initialization only")
	flag.Parse()

	gConfig = config.NewConfig("hls.cfg")
	gDB = database.NewKeyValueDB("hls.db")

	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	scheduler := cron.New(cron.WithLocation(jakartaTime))
	scheduler.AddFunc(gConfig.CleanCron(), func() {
		gDB.Clean()
	})

	sett := settings.NewSettings()
	privKey := getPrivKey(
		sett,
		"hls.stg",
		utils.InputString("Storage password: "),
		utils.InputString("Object password: "),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gClient = local.NewClient(privKey, sett)
	if gClient == nil {
		return fmt.Errorf("failed create client node")
	}

	fmt.Printf("Public key: %s\n", privKey.PubKey())
	if initOnly {
		os.Exit(0)
	}

	fmt.Printf("Service is listening [%s]...\n", gConfig.Address())
	return nil
}

func getPrivKey(sett settings.ISettings, filepath, storagePasw, objectPasw string) crypto.IPrivKey {
	storage := local.NewStorage(
		sett,
		filepath,
		storagePasw,
	)
	if storage == nil {
		return nil
	}

	bpriv, err := storage.Read("private_key", objectPasw)
	if err == nil {
		return crypto.LoadPrivKey(bpriv)
	}

	priv := crypto.NewPrivKey(cAKeySize)
	err = storage.Write("private_key", objectPasw, priv.Bytes())
	if err != nil {
		return nil
	}

	return priv
}
