package main

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/cmd/hmc/config"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/storage"
	"github.com/number571/go-peer/utils"

	hmc_settings "github.com/number571/go-peer/cmd/hmc/settings"
	hms_database "github.com/number571/go-peer/cmd/hms/database"
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

func hmcDefaultInit() error {
	gWrapper = newWrapper(
		newConfigWrapper(config.NewConfig("hmc.cfg")),
	)
	gDB = hms_database.NewKeyValueDB("hmc.db")

	storageSettings := hmc_settings.NewSettings()
	privKey := getPrivKey(
		storageSettings,
		"hmc.stg",
		[]byte(utils.NewInput(storageSettings, "Password#Stg: ").Password()),
		[]byte(utils.NewInput(storageSettings, "Password#Obj: ").Password()),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gActions = newActions()
	gClient = client.NewClient(hms_settings.NewSettings(), privKey)

	return nil
}

func getPrivKey(sett settings.ISettings, filepath string, storageKey, objectKey []byte) asymmetric.IPrivKey {
	// create/open storage
	storage := storage.NewCryptoStorage(
		sett,
		filepath,
		storageKey,
	)
	if storage == nil {
		return nil
	}

	// get private key
	bpriv, err := storage.Get(objectKey)
	if err == nil {
		return asymmetric.LoadRSAPrivKey(bpriv)
	}

	// private key not exist
	answ := utils.NewInput(nil, "Private key by password not exist.\nGenerate new? [y/n]: ").String()
	switch strings.ToLower(answ) {
	case "y", "yes":
		// generate private key
		priv := asymmetric.NewRSAPrivKey(cAKeySize)
		err := storage.Set(objectKey, priv.Bytes())
		if err != nil {
			panic(err)
		}
		return priv
	case "n", "no":
		// exit from program
		return nil
	default:
		// undefined answer
		panic("input answer not equal [y/yes, n/no]")
	}
}
