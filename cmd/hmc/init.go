package main

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/cmd/hmc/config"
	"github.com/number571/go-peer/cmd/hmc/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage"

	hms_database "github.com/number571/go-peer/cmd/hms/database"
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

func hmcDefaultInit() error {
	gWrapper = newWrapper(
		newConfigWrapper(config.NewConfig("hmc.cfg")),
	)
	gDB = hms_database.NewKeyValueDB("hmc.db")

	privKey := getPrivKey(
		"hmc.stg",
		[]byte(newInputter("Password#Stg: ").Password()),
		[]byte(newInputter("Password#Obj: ").Password()),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gActions = newActions()
	gClient = client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    hms_settings.CSizeWork,
			FMessageSize: hms_settings.CSizePack,
		}),
		privKey,
	)

	return nil
}

func getPrivKey(filepath string, storageKey, objectKey []byte) asymmetric.IPrivKey {
	// create/open storage
	storage, err := storage.NewCryptoStorage(
		filepath,
		storageKey,
		settings.CWorkSize,
	)
	if err != nil {
		return nil
	}

	// get private key
	bpriv, err := storage.Get(objectKey)
	if err == nil {
		return asymmetric.LoadRSAPrivKey(bpriv)
	}

	// private key not exist
	answ := newInputter("Private key by password not exist.\nGenerate new? [y/n]: ").String()
	switch strings.ToLower(answ) {
	case "y", "yes":
		// generate private key
		priv := asymmetric.NewRSAPrivKey(settings.CAKeySize)
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
