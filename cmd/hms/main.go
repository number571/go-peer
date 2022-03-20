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
	response(w, 0, []byte("hidden message service"))
}

func sizePage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Receiver []byte `json:"receiver"`
	}

	if r.Method != "POST" {
		response(w, cErrorMethod, []byte("failed: method POST"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response(w, cErrorDecode, []byte("failed: decode request"))
		return
	}

	if len(request.Receiver) != crypto.HashSize {
		response(w, cErrorSize, []byte("failed: receiver size"))
		return
	}

	size := gDB.Size(request.Receiver)
	response(w, cErrorNone, encoding.Uint64ToBytes(size))
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Receiver []byte `json:"receiver"`
		Index    uint64 `json:"index"`
	}

	if r.Method != "POST" {
		response(w, cErrorMethod, []byte("failed: method POST"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response(w, cErrorDecode, []byte("failed: decode request"))
		return
	}

	if len(request.Receiver) != crypto.HashSize {
		response(w, cErrorSize, []byte("failed: receiver size"))
		return
	}

	msg := gDB.Load(request.Receiver, request.Index)
	if msg == nil {
		response(w, cErrorLoad, []byte("failed: load message"))
		return
	}

	response(w, cErrorNone, msg.ToPackage().Bytes())
}

func pushPage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Receiver []byte `json:"receiver"`
		Package  []byte `json:"package"`
	}

	if r.Method != "POST" {
		response(w, cErrorMethod, []byte("failed: method POST"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response(w, cErrorDecode, []byte("failed: decode request"))
		return
	}

	if len(request.Receiver) != crypto.HashSize {
		response(w, cErrorSize, []byte("failed: receiver size"))
		return
	}

	msg := local.LoadPackage(request.Package).ToMessage()
	if msg == nil {
		response(w, cErrorMessage, []byte("failed: decode message"))
		return
	}

	err = gDB.Push(request.Receiver, msg)
	if err != nil {
		response(w, cErrorPush, []byte("failed: push message"))
		return
	}

	response(w, cErrorNone, []byte("success"))
}

func response(w http.ResponseWriter, ret int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Result []byte `json:"result"`
		Return int    `json:"return"`
	}{
		Result: res,
		Return: ret,
	})
}
