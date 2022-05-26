// HMS - Hidden Message Service
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

func main() {
	err := hmsDefaultInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/size", sizePage)
	http.HandleFunc("/load", loadPage)
	http.HandleFunc("/push", pushPage)

	err = http.ListenAndServe(gConfig.Address(), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	response(w, hms_settings.CErrorNone, []byte("hidden message service"))
}

func sizePage(w http.ResponseWriter, r *http.Request) {
	var vRequest hms_settings.SSizeRequest

	if r.Method != "POST" {
		response(w, hms_settings.CErrorMethod, []byte("failed: incorrect method"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		response(w, hms_settings.CErrorDecode, []byte("failed: decode request"))
		return
	}

	size, err := gDB.Size(vRequest.Receiver)
	if err != nil {
		response(w, hms_settings.CErrorLoad, []byte("failed: load size"))
		return
	}

	response(w, hms_settings.CErrorNone, encoding.Uint64ToBytes(size))
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	var vRequest hms_settings.SLoadRequest

	if r.Method != "POST" {
		response(w, hms_settings.CErrorMethod, []byte("failed: incorrect method"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		response(w, hms_settings.CErrorDecode, []byte("failed: decode request"))
		return
	}

	msg, err := gDB.Load(vRequest.Receiver, vRequest.Index)
	if err != nil {
		response(w, hms_settings.CErrorLoad, []byte("failed: load message"))
		return
	}

	response(w, hms_settings.CErrorNone, msg.ToPackage().Bytes())
}

func pushPage(w http.ResponseWriter, r *http.Request) {
	var vRequest hms_settings.SPushRequest

	if r.Method != "POST" {
		response(w, hms_settings.CErrorMethod, []byte("failed: incorrect method"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		response(w, hms_settings.CErrorDecode, []byte("failed: decode request"))
		return
	}

	msg := local.LoadPackage(vRequest.Package).ToMessage()
	if msg == nil {
		response(w, hms_settings.CErrorMessage, []byte("failed: decode message"))
		return
	}

	puzzle := crypto.NewPuzzle(gSettings.Get(settings.SizeWork))
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		response(w, hms_settings.CErrorWorkSize, []byte("failed: incorrect work size"))
		return
	}

	err = gDB.Push(vRequest.Receiver, msg)
	if err != nil {
		response(w, hms_settings.CErrorPush, []byte("failed: push message"))
		return
	}

	response(w, hms_settings.CErrorNone, []byte("success"))
}

func response(w http.ResponseWriter, ret int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hms_settings.SResponse{
		Result: res,
		Return: ret,
	})
}
