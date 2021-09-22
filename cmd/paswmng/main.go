// go run main.go -sn storage.enc -sp st-password --set -n example.com password
// go run main.go -sn storage.enc -sp st-password --get -n example.com
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/number571/gopeer/local"
)

var (
	STORAGE_NAME     string
	STORAGE_PASSWORD string

	SUBJECT_NAME     string
	SUBJECT_PASSWORD string

	IS_SET_SUBJECT bool
	IS_GET_SUBJECT bool

	MOST_SECRET string
)

func init() {
	const (
		storageName     = "load storage by filename"
		storagePassword = "storage password for decrypt main structure"

		subjectName     = "load subject"
		subjectPassword = "load password"

		storageSet = "option for set subject and password"
		storageGet = "option for get password from subject"

		mostSecret = "set key for most secret password"
	)

	flag.StringVar(&STORAGE_NAME, "sn", "", storageName)
	flag.StringVar(&STORAGE_PASSWORD, "sp", "", storagePassword)

	flag.StringVar(&SUBJECT_NAME, "n", "", subjectName)
	flag.StringVar(&SUBJECT_PASSWORD, "p", "", subjectPassword)

	flag.BoolVar(&IS_SET_SUBJECT, "set", false, storageSet)
	flag.BoolVar(&IS_GET_SUBJECT, "get", false, storageGet)

	flag.StringVar(&MOST_SECRET, "ms", "null", mostSecret)

	flag.Parse()

	switch {
	case IS_GET_SUBJECT && IS_SET_SUBJECT:
		fmt.Println("error: get and set in one request")
		os.Exit(1)
	case STORAGE_NAME == "" || STORAGE_PASSWORD == "":
		fmt.Println("error: storage filename or password is null")
		os.Exit(2)
	case SUBJECT_NAME == "":
		fmt.Println("error: subject name is null")
		os.Exit(3)
	case IS_SET_SUBJECT && SUBJECT_PASSWORD == "":
		fmt.Println("error: subject name or password is null with set option")
		os.Exit(4)
	}
}

func main() {
	var (
		defaultKey = MOST_SECRET
		store      = local.NewStorage(STORAGE_NAME, STORAGE_PASSWORD)
	)

	switch {
	case IS_GET_SUBJECT:
		key, err := store.Read(SUBJECT_NAME, defaultKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(10)
		}
		fmt.Println(string(key))
	case IS_SET_SUBJECT:
		err := store.Write([]byte(SUBJECT_PASSWORD), SUBJECT_NAME, defaultKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(11)
		}
	}
}
