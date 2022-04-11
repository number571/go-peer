package main

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/cmd/hms/utils"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
)

func pageIndex(w http.ResponseWriter, r *http.Request) {
	response(w, cErrorNone, []byte("hidden lake service"))
}

func pageStatus(w http.ResponseWriter, r *http.Request) {
	response(w, cErrorNone, utils.Serialize(struct {
		PubKey []byte `json:"pubkey"`
	}{
		PubKey: gNode.Client().PubKey().Bytes(),
	}))
}

func pageMessage(w http.ResponseWriter, r *http.Request) {
	var vRequest struct {
		Receiver []byte `json:"receiver"`
		Title    []byte `json:"title"`
		Data     []byte `json:"data"`
	}

	if r.Method != "POST" {
		response(w, cErrorMethod, []byte("failed: method POST"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		response(w, cErrorDecodeRequest, []byte("failed: decode request"))
		return
	}

	pubKey := crypto.LoadPubKey(vRequest.Receiver)
	if pubKey == nil {
		response(w, cErrorDecodePubKey, []byte("failed: decode public key"))
		return
	}

	inOnline := []crypto.IPubKey{}
	for _, val := range gNode.Checker().ListWithInfo() {
		if !val.Online() {
			continue
		}
		inOnline = append(inOnline, val.PubKey())
	}

	rand := crypto.NewPRNG()
	randSizeRoute := rand.Uint64() % cSizeRoute

	resp, err := gNode.Request(
		local.NewRoute(pubKey).
			WithRedirects(
				gPPrivKey,
				local.NewSelector(inOnline).
					Shuffle().
					Return(randSizeRoute),
			),
		local.NewMessage(vRequest.Title, vRequest.Data),
	)
	if err != nil {
		response(w, cErrorResponseMessage, []byte("failed: response message"))
		return
	}

	response(w, cErrorNone, resp)
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
