package main

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hms/hmc"
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/crypto/puzzle"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/offline/message"
	"github.com/number571/go-peer/settings"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	response(w, hms_settings.CErrorNone, []byte(hms_settings.CTitlePattern))
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

	if uint64(len(vRequest.Package)) > gSettings.Get(settings.CSizePack) {
		response(w, hms_settings.CErrorPackSize, []byte("failed: incorrect package size"))
		return
	}

	msg := message.LoadPackage(vRequest.Package).ToMessage()
	if msg == nil {
		response(w, hms_settings.CErrorMessage, []byte("failed: decode message"))
		return
	}

	puzzle := puzzle.NewPoWPuzzle(gSettings.Get(settings.CSizeWork))
	if !puzzle.Verify(msg.Body().Hash(), msg.Body().Proof()) {
		response(w, hms_settings.CErrorWorkSize, []byte("failed: incorrect work size"))
		return
	}

	err = gDB.Push(vRequest.Receiver, msg)
	if err != nil {
		response(w, hms_settings.CErrorPush, []byte("failed: push message"))
		return
	}

	go func() {
		for _, host := range gConfig.Connections() {
			hmc.NewRequester(host).Push(&vRequest)
		}
	}()

	response(w, hms_settings.CErrorNone, []byte("success"))
}

func response(w http.ResponseWriter, ret int, res []byte) {
	w.Header().Set("Content-Type", hms_settings.CContentType)
	json.NewEncoder(w).Encode(&hms_settings.SResponse{
		Result: res,
		Return: ret,
	})
}
