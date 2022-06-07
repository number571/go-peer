// HMS - Hidden Message Service
package main

import (
	"fmt"
	"net/http"
	"os"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

func main() {
	err := hmsDefaultInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc("/", indexPage)
	http.HandleFunc(hms_settings.CSizePath, sizePage)
	http.HandleFunc(hms_settings.CLoadPath, loadPage)
	http.HandleFunc(hms_settings.CPushSize, pushPage)

	err = http.ListenAndServe(gConfig.Address(), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
