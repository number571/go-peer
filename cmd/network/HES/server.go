package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cr "github.com/number571/gopeer/crypto"
	en "github.com/number571/gopeer/encoding"
	lc "github.com/number571/gopeer/local"
)

var (
	DATABASE = NewDB("server.db")
	FLCONFIG = NewCFG("server.cfg")
)

func init() {
	go delOldEmailsByTime(24*time.Hour, 6*time.Hour)
	hesDefaultInit("localhost:8080")
	fmt.Printf("Server is listening [%s] ...\n\n", OPENADDR)
}

func delOldEmailsByTime(deltime, period time.Duration) {
	for {
		DATABASE.DelEmailsByTime(deltime)
		time.Sleep(period)
	}
}

func main() {
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/email/send", emailSendPage)
	http.HandleFunc("/email/recv", emailRecvPage)
	http.ListenAndServe(OPENADDR, nil)
}

func help() string {
	return `
1. exit   - close server;
2. help   - commands info;
3. list   - list connections;
4. append - append connect to list;
5. delete - delete connect from list;
`
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Macp string `json:"macp"`
	}
	if r.Method != "POST" {
		response(w, 0, "hidden email service")
		return
	}
	if r.ContentLength > int64(MAXESIZE) {
		response(w, 1, "error: max size")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, 2, "error: parse json")
		return
	}
	pasw := cr.NewSHA256([]byte(FLCONFIG.Pasw)).Bytes()
	cipher := cr.NewCipher(pasw)
	dect := cipher.Decrypt(en.Base64Decode(req.Macp))
	if !bytes.Equal([]byte(TMESSAGE), dect) {
		response(w, 3, "error: message authentication code")
		return
	}
	response(w, 0, "success: check connection")
}

func emailSendPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Recv string `json:"recv"`
		Data string `json:"data"`
		Macp string `json:"macp"`
	}
	if r.Method != "POST" {
		response(w, 1, "error: method != POST")
		return
	}
	if r.ContentLength > int64(MAXESIZE) {
		response(w, 2, "error: max size")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, 3, "error: parse json")
		return
	}
	pack := lc.Package(req.Data).Deserialize()
	if pack == nil {
		response(w, 4, "error: deserialize package")
		return
	}
	hash := pack.Body.Hash
	puzzle := cr.NewPuzzle(uint8(POWSDIFF))
	if !puzzle.Verify(hash, pack.Body.Npow) {
		response(w, 5, "error: proof of work")
		return
	}
	pasw := cr.NewSHA256([]byte(FLCONFIG.Pasw)).Bytes()
	cipher := cr.NewCipher(pasw)
	dech := cipher.Decrypt(en.Base64Decode(req.Macp))
	if !bytes.Equal(hash, dech) {
		response(w, 6, "error: message authentication code")
		return
	}
	err = DATABASE.SetEmail(req.Recv, pack)
	if err != nil {
		response(w, 7, "error: save email")
		return
	}
	for _, conn := range FLCONFIG.Conns {
		go func() {
			addr := conn[0]
			pasw := cr.NewSHA256([]byte(conn[1])).Bytes()
			cipher := cr.NewCipher(pasw)
			req.Macp = en.Base64Encode(cipher.Encrypt(hash))
			resp, err := HTCLIENT.Post(
				addr+"/email/send",
				"application/json",
				bytes.NewReader(serialize(req)),
			)
			if err != nil {
				return
			}
			resp.Body.Close()
		}()
	}
	response(w, 0, "success: email saved")
}

func emailRecvPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Recv string `json:"recv"`
		Data int    `json:"data"`
	}
	if r.Method != "POST" {
		response(w, 1, "error: method != POST")
		return
	}
	if r.ContentLength > int64(MAXESIZE) {
		response(w, 2, "error: max size")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, 3, "error: parse json")
		return
	}
	if req.Data == 0 {
		response(w, 0, fmt.Sprintf("%d", DATABASE.Size(req.Recv)))
		return
	}
	res := DATABASE.GetEmail(req.Data, req.Recv)
	if res == "" {
		response(w, 4, "error: nothing data")
		return
	}
	response(w, 0, res)
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", "application/json")
	var resp struct {
		Result string `json:"result"`
		Return int    `json:"return"`
	}
	resp.Result = res
	resp.Return = ret
	json.NewEncoder(w).Encode(resp)
}
