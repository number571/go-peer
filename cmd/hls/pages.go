package main

import (
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
)

type sDefaultResponse struct {
	Result string `json:"result"`
	Return int    `json:"return"`
}

func pageIndex(w http.ResponseWriter, r *http.Request) {
	defaultHeaders(w)
	defaultResponse(w, cErrorNone, "hidden lake service")
}

func pageStatus(w http.ResponseWriter, r *http.Request) {
	defaultHeaders(w)

	type sStatus struct {
		PubKey string `json:"pub_key"`
		Online bool   `json:"online"`
	}

	type sCustomResponse struct {
		PubKey  string     `json:"pub_key"`
		Network []*sStatus `json:"network"`
		sDefaultResponse
	}

	if r.Method != "GET" {
		defaultResponse(w, cErrorMethod, "failed: incorrect method")
		return
	}

	var network []*sStatus
	for _, info := range gNode.Checker().ListWithInfo() {
		network = append(network, &sStatus{
			PubKey: info.PubKey().String(),
			Online: info.Online(),
		})
	}

	json.NewEncoder(w).Encode(sCustomResponse{
		PubKey:  gNode.Client().PubKey().String(),
		Network: network,
		sDefaultResponse: sDefaultResponse{
			Result: "success",
			Return: cErrorNone,
		},
	})
}

func pageMessage(w http.ResponseWriter, r *http.Request) {
	defaultHeaders(w)

	var vRequest struct {
		Receiver string `json:"receiver"` // public key
		Title    []byte `json:"title"`
		Data     []byte `json:"data"`
	}

	if r.Method != "POST" {
		defaultResponse(w, cErrorMethod, "failed: incorrect method")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&vRequest)
	if err != nil {
		defaultResponse(w, cErrorDecodeRequest, "failed: decode request")
		return
	}

	pubKey := crypto.LoadPubKey(vRequest.Receiver)
	if pubKey == nil {
		defaultResponse(w, cErrorDecodePubKey, "failed: decode public key")
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
		defaultResponse(w, cErrorResponseMessage, "failed: response message")
		return
	}

	defaultResponse(w, cErrorNone, encoding.Base64Encode(resp))
}

func defaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func defaultResponse(w http.ResponseWriter, ret int, res string) {
	json.NewEncoder(w).Encode(sDefaultResponse{
		Result: res,
		Return: ret,
	})
}
