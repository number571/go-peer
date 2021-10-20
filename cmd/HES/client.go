package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	gp "github.com/number571/gopeer"
	cr "github.com/number571/gopeer/crypto"
	en "github.com/number571/gopeer/encoding"
	lc "github.com/number571/gopeer/local"
)

type TemplateResult struct {
	Auth   string
	Result string
	Return int
}

const (
	FSEPARAT = "\001\007\005\000\005\007\001"
	MAXEPAGE = 5 // view emails in one page
	MAXCOUNT = 5 // load emails from one node
)

const (
	RET_SUCCESS = 0
	RET_DANGER  = 1
	RET_WARNING = 2
)

const (
	PATH_VIEWS  = "userside/views/"
	PATH_STATIC = "userside/static/"
)

var (
	DATABASE = NewDB("client.db")
	SESSIONS = NewSessions()
)

func init() {
	go delOldSessionsByTime(1*time.Hour, 15*time.Minute)
	hesDefaultInit("localhost:7545")
	fmt.Printf("Client is listening [%s] ...\n\n", OPENADDR)
}

func delOldSessionsByTime(deltime, period time.Duration) {
	for {
		SESSIONS.DelByTime(deltime)
		time.Sleep(period)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(http.Dir(PATH_STATIC))),
	)
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/account", accountPage)
	http.HandleFunc("/account/public_key", accountPublicKeyPage)
	http.HandleFunc("/account/private_key", accountPrivateKeyPage)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signin", signinPage)
	http.HandleFunc("/signout", signoutPage)
	http.HandleFunc("/network", networkPage)
	http.HandleFunc("/network/read", networkReadPage)
	http.HandleFunc("/network/write", networkWritePage)
	http.HandleFunc("/network/contact", networkContactPage)
	http.HandleFunc("/network/connect", networkConnectPage)
	http.ListenAndServe(OPENADDR, nil)
}

func handleFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			indexPage(w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"index.html",
	)
	if err != nil {
		panic("error: load index.html")
	}
	t.Execute(w, TemplateResult{
		Auth: getName(SESSIONS.Get(r)),
	})
}

func signupPage(w http.ResponseWriter, r *http.Request) {
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"signup.html",
	)
	if err != nil {
		panic("error: load signup.html")
	}
	if SESSIONS.Get(r) != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" {
		name := r.FormValue("username")
		pasw := r.FormValue("password")
		spriv := r.FormValue("private_key")
		priv := cr.LoadPrivKeyByString(spriv)
		if pasw != r.FormValue("password_repeat") {
			retcod, result = makeResult(RET_DANGER, "error: passwords not equal")
			goto close
		}
		if spriv != "" && priv == nil {
			retcod, result = makeResult(RET_DANGER, "error: private key is not valid")
			goto close
		}
		if priv == nil {
			priv = cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))
		}
		err := DATABASE.SetUser(name, pasw, priv)
		if err != nil {
			retcod, result = makeResult(RET_DANGER,
				fmt.Sprintf("error: %s", err.Error()))
			goto close
		}
		http.Redirect(w, r, "/signin", 302)
		return
	}
close:
	t.Execute(w, TemplateResult{
		Result: result,
		Return: retcod,
	})
}

func signinPage(w http.ResponseWriter, r *http.Request) {
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"signin.html",
	)
	if err != nil {
		panic("error: load signin.html")
	}
	if SESSIONS.Get(r) != nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" {
		name := r.FormValue("username")
		pasw := r.FormValue("password")
		user := DATABASE.GetUser(name, pasw)
		if user == nil {
			retcod, result = makeResult(RET_DANGER, "error: username of password incorrect")
			goto close
		}
		SESSIONS.Set(w, user)
		http.Redirect(w, r, "/", 302)
		return
	}
close:
	t.Execute(w, TemplateResult{
		Result: result,
		Return: retcod,
	})
}

func signoutPage(w http.ResponseWriter, r *http.Request) {
	SESSIONS.Del(w, r)
	http.Redirect(w, r, "/", 302)
}

func accountPage(w http.ResponseWriter, r *http.Request) {
	type AccountTemplateResult struct {
		TemplateResult
		PublicKey  string
		PrivateKey string
	}
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"account.html",
	)
	if err != nil {
		panic("error: load account.html")
	}
	user := SESSIONS.Get(r)
	if user == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" && r.FormValue("delete") != "" {
		name := r.FormValue("username")
		pasw := r.FormValue("password")
		cuser := DATABASE.GetUser(name, pasw)
		if cuser == nil || cuser.Id != user.Id {
			retcod, result = makeResult(RET_DANGER, "error: username of password incorrect")
			goto close
		}
		SESSIONS.Del(w, r)
		DATABASE.DelUser(cuser)
		http.Redirect(w, r, "/", 302)
		return
	}
close:
	t.Execute(w, AccountTemplateResult{
		TemplateResult: TemplateResult{
			Auth:   getName(SESSIONS.Get(r)),
			Result: result,
			Return: retcod,
		},
		PublicKey:  user.Priv.PubKey().String(),
		PrivateKey: user.Priv.String(),
	})
}

func accountPublicKeyPage(w http.ResponseWriter, r *http.Request) {
	user := SESSIONS.Get(r)
	if user == nil {
		fmt.Fprint(w, "error: session is null")
		return
	}
	dataString := user.Priv.PubKey().String()
	qrCode, err := qr.Encode(dataString, qr.Q, qr.Auto)
	if err != nil {
		fmt.Fprint(w, "error: qrcode generate")
		return
	}
	qrCode, err = barcode.Scale(qrCode, 768, 768)
	if err != nil {
		fmt.Fprint(w, "error: qrcode scale")
		return
	}
	png.Encode(w, qrCode)
}

func accountPrivateKeyPage(w http.ResponseWriter, r *http.Request) {
	user := SESSIONS.Get(r)
	if user == nil {
		fmt.Fprint(w, "error: session is null")
		return
	}
	dataString := user.Priv.String()
	qrCode, err := qr.Encode(dataString, qr.L, qr.Auto)
	if err != nil {
		fmt.Fprint(w, "error: qrcode generate")
		return
	}
	qrCode, err = barcode.Scale(qrCode, 1024, 1024)
	if err != nil {
		fmt.Fprint(w, "error: qrcode scale")
		return
	}
	png.Encode(w, qrCode)
}

func networkPage(w http.ResponseWriter, r *http.Request) {
	type ReadTemplateResult struct {
		TemplateResult
		Page   int
		Emails []Email
	}
	page := 0
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.New("base.html").Funcs(template.FuncMap{
		"inc":   func(x int) int { return x + 1 },
		"dec":   func(x int) int { return x - 1 },
		"texts": getTexts,
	}).ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"network.html",
	)
	if err != nil {
		panic("error: load network.html")
	}
	t = template.Must(t, err)
	user := SESSIONS.Get(r)
	if user == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "GET" && r.FormValue("page") != "" {
		num, err := strconv.Atoi(r.FormValue("page"))
		if err != nil {
			retcod, result = makeResult(RET_DANGER, "error: parse atoi")
			goto close
		}
		page = num
	}
	if r.Method == "POST" && r.FormValue("delete") != "" {
		hash := r.FormValue("email")
		DATABASE.DelEmail(user, hash)
	}
	if r.Method == "POST" && r.FormValue("update") != "" {
		conns := DATABASE.GetConns(user)
		for _, conn := range conns {
			go readEmails(user, conn[0])
		}
		time.Sleep(3 * time.Second)
	}
close:
	t.Execute(w, ReadTemplateResult{
		TemplateResult: TemplateResult{
			Auth:   getName(SESSIONS.Get(r)),
			Result: result,
			Return: retcod,
		},
		Page:   page,
		Emails: DATABASE.GetEmails(user, page*MAXEPAGE, MAXEPAGE),
	})
}

// FS = FSEPARAT
// head = title   || FS || filename[0]     || ... || FS || filename[n]
// body = message || FS || base64(file[0]) || ... || FS || base64(file[n])
func networkWritePage(w http.ResponseWriter, r *http.Request) {
	type WriteTemplateResult struct {
		TemplateResult
		Contacts map[string]string
	}
	type Resp struct {
		Result string `json:"result"`
		Return int    `json:"return"`
	}
	type Req struct {
		Recv string `json:"recv"`
		Data string `json:"data"`
		Macp string `json:"macp"`
	}
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"write.html",
	)
	if err != nil {
		panic("error: load write.html")
	}
	user := SESSIONS.Get(r)
	if user == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" {
		err := r.ParseMultipartForm(int64(MAXESIZE))
		if err != nil {
			retcod, result = makeResult(RET_DANGER, "error: max size")
			goto close
		}
		recv := cr.LoadPubKeyByString(r.FormValue("receiver"))
		if recv == nil {
			retcod, result = makeResult(RET_DANGER, "error: receiver is null")
			goto close
		}
		head := strings.TrimSpace(r.FormValue("title"))
		body := strings.TrimSpace(r.FormValue("message"))
		if head == "" || body == "" {
			retcod, result = makeResult(RET_DANGER, "error: head or body is null")
			goto close
		}
		files := r.MultipartForm.File["files"]
		for i := range files {
			file, err := files[i].Open()
			if err != nil {
				retcod, result = makeResult(RET_DANGER, "error: open file")
				goto close
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				retcod, result = makeResult(RET_DANGER, "error: read file")
				goto close
			}
			file.Close()
			head += FSEPARAT + files[i].Filename
			body += FSEPARAT + en.Base64Encode(content)
		}
		client := lc.NewClient(user.Priv)
		pack := client.Encrypt(recv, newEmail(user.Name, head, body))
		hash := pack.Body.Hash
		conns := DATABASE.GetConns(user)
		req := Req{
			Recv: string(recv.Address()),
			Data: string(pack.Serialize()),
		}
		if uint(len(req.Data)) > MAXESIZE {
			retcod, result = makeResult(RET_DANGER, "error: max size")
			goto close
		}
		for _, conn := range conns {
			pasw := cr.SumHash([]byte(conn[1]))
			cipher := cr.NewCipher(pasw)
			req.Macp = en.Base64Encode(cipher.Encrypt(hash))
			go writeEmails(conn[0], serialize(req))
		}
		result = "success: email send"
	}
close:
	t.Execute(w, WriteTemplateResult{
		TemplateResult: TemplateResult{
			Auth:   getName(SESSIONS.Get(r)),
			Result: result,
			Return: retcod,
		},
		Contacts: DATABASE.GetContacts(user),
	})
}

func networkReadPage(w http.ResponseWriter, r *http.Request) {
	type ReadTemplateResult struct {
		TemplateResult
		Email *Email
	}
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.New("base.html").Funcs(template.FuncMap{
		"split": strings.Split,
		"texts": getTexts,
		"files": getFiles,
	}).ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"read.html",
	)
	if err != nil {
		panic("error: load read.html")
	}
	t = template.Must(t, err)
	user := SESSIONS.Get(r)
	if user == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" {
		pub := cr.LoadPubKeyByString(r.FormValue("public_key"))
		if pub == nil {
			fmt.Fprint(w, "error: public key is null")
			return
		}
		dataString := pub.String()
		qrCode, err := qr.Encode(dataString, qr.Q, qr.Auto)
		if err != nil {
			fmt.Fprint(w, "error: qrcode generate")
			return
		}
		qrCode, err = barcode.Scale(qrCode, 768, 768)
		if err != nil {
			fmt.Fprint(w, "error: qrcode scale")
			return
		}
		png.Encode(w, qrCode)
		return
	}
	var email *Email
	id, err := strconv.Atoi(r.FormValue("email"))
	if err != nil {
		retcod, result = makeResult(RET_DANGER, "error: atoi parse")
		goto close
	}
	email = DATABASE.GetEmail(user, id)
	if email == nil {
		retcod, result = makeResult(RET_DANGER, "error: email undefined")
		goto close
	}
close:
	t.Execute(w, ReadTemplateResult{
		TemplateResult: TemplateResult{
			Auth:   getName(SESSIONS.Get(r)),
			Result: result,
			Return: retcod,
		},
		Email: email,
	})
}

func networkContactPage(w http.ResponseWriter, r *http.Request) {
	type ContactTemplateResult struct {
		TemplateResult
		F2F      bool
		Contacts map[string]string
	}
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"contact.html",
	)
	if err != nil {
		panic("error: load contact.html")
	}
	user := SESSIONS.Get(r)
	if user == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" && r.FormValue("switchf2f") != "" {
		DATABASE.SwitchF2F(user)
	}
	if r.Method == "POST" && r.FormValue("append") != "" {
		name := r.FormValue("nickname")
		publ := cr.LoadPubKeyByString(r.FormValue("public_key"))
		err := DATABASE.SetContact(user, name, publ)
		if err != nil {
			retcod, result = makeResult(RET_DANGER,
				fmt.Sprintf("error: %s", err.Error()))
			goto close
		}
	}
	if r.Method == "POST" && r.FormValue("delete") != "" {
		publ := cr.LoadPubKeyByString(r.FormValue("public_key"))
		err := DATABASE.DelContact(user, publ)
		if err != nil {
			retcod, result = makeResult(RET_DANGER,
				fmt.Sprintf("error: %s", err.Error()))
			goto close
		}
	}
close:
	t.Execute(w, ContactTemplateResult{
		TemplateResult: TemplateResult{
			Auth:   getName(SESSIONS.Get(r)),
			Result: result,
			Return: retcod,
		},
		F2F:      DATABASE.StateF2F(user),
		Contacts: DATABASE.GetContacts(user),
	})
}

func networkConnectPage(w http.ResponseWriter, r *http.Request) {
	type ConnTemplateResult struct {
		TemplateResult
		Connects [][2]string
	}
	retcod, result := makeResult(RET_SUCCESS, "")
	t, err := template.ParseFiles(
		PATH_VIEWS+"base.html",
		PATH_VIEWS+"connect.html",
	)
	if err != nil {
		panic("error: load connect.html")
	}
	user := SESSIONS.Get(r)
	if user == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	if r.Method == "POST" && r.FormValue("check") != "" {
		conns := DATABASE.GetConns(user)
		for _, conn := range conns {
			ret, res := checkConnection(conn)
			if ret != RET_SUCCESS {
				result += res
				retcod = RET_WARNING
			}
		}
		if retcod != RET_SUCCESS {
			goto close
		}
		result = "success: all connections work"
	}
	if r.Method == "POST" && r.FormValue("append") != "" {
		host := r.FormValue("hostname")
		pasw := r.FormValue("password")
		err := DATABASE.SetConn(user, host, pasw)
		if err != nil {
			retcod, result = makeResult(RET_DANGER,
				fmt.Sprintf("error: %s", err.Error()))
			goto close
		}
	}
	if r.Method == "POST" && r.FormValue("delete") != "" {
		host := r.FormValue("hostname")
		err := DATABASE.DelConn(user, host)
		if err != nil {
			retcod, result = makeResult(RET_DANGER,
				fmt.Sprintf("error: %s", err.Error()))
			goto close
		}
	}
close:
	t.Execute(w, ConnTemplateResult{
		TemplateResult: TemplateResult{
			Auth:   getName(SESSIONS.Get(r)),
			Result: result,
			Return: retcod,
		},
		Connects: DATABASE.GetConns(user),
	})
}

func checkConnection(conn [2]string) (int, string) {
	type Resp struct {
		Result string `json:"result"`
		Return int    `json:"return"`
	}
	type Req struct {
		Macp string `json:"macp"`
	}
	var servresp Resp
	pasw := cr.SumHash([]byte(conn[1]))
	cipher := cr.NewCipher(pasw)
	macp := cipher.Encrypt([]byte(TMESSAGE))
	resp, err := HTCLIENT.Post(
		strings.TrimRight(conn[0], " /")+"/",
		"application/json",
		bytes.NewReader(serialize(Req{
			Macp: en.Base64Encode(macp),
		})),
	)
	if err != nil {
		return makeResult(RET_DANGER,
			fmt.Sprintf("%s='%s';\n", conn[0], "error: connect"))
	}
	if resp.ContentLength > int64(MAXESIZE) {
		return makeResult(RET_DANGER,
			fmt.Sprintf("%s='%s';\n", conn[0], "error: max size"))
	}
	err = json.NewDecoder(resp.Body).Decode(&servresp)
	resp.Body.Close()
	if err != nil {
		return makeResult(RET_DANGER,
			fmt.Sprintf("%s='%s';\n", conn[0], "error: parse json"))
	}
	if servresp.Return != 0 {
		return makeResult(RET_DANGER,
			fmt.Sprintf("%s='%s';\n", conn[0], servresp.Result))
	}
	return makeResult(RET_SUCCESS, "")
}

func writeEmails(addr string, rdata []byte) {
	type Resp struct {
		Result string `json:"result"`
		Return int    `json:"return"`
	}
	type Req struct {
		Recv string `json:"recv"`
		Data int    `json:"data"`
	}
	var servresp Resp
	resp, err := HTCLIENT.Post(
		strings.TrimRight(addr, " /")+"/email/send",
		"application/json",
		bytes.NewReader(rdata),
	)
	if err != nil {
		return
	}
	if resp.ContentLength > int64(MAXESIZE) {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&servresp)
	resp.Body.Close()
	if err != nil {
		return
	}
	if servresp.Return != 0 {
		return
	}
}

func readEmails(user *User, addr string) {
	type Resp struct {
		Result string `json:"result"`
		Return int    `json:"return"`
	}
	type Req struct {
		Recv string `json:"recv"`
		Data int    `json:"data"`
	}
	var servresp Resp
	client := lc.NewClient(user.Priv)
	pbhash := string(client.PubKey().Address())
	// GET SIZE EMAILS
	resp, err := HTCLIENT.Post(
		strings.TrimRight(addr, " /")+"/email/recv",
		"application/json",
		bytes.NewReader(serialize(Req{
			Recv: pbhash,
			Data: 0,
		})),
	)
	if err != nil {
		return
	}
	if resp.ContentLength > int64(MAXESIZE) {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&servresp)
	resp.Body.Close()
	if err != nil {
		return
	}
	if servresp.Return != 0 {
		return
	}
	// GET DATA EMAILS
	size, err := strconv.Atoi(servresp.Result)
	if err != nil {
		return
	}
	for i, count := 1, 0; i <= size; i++ {
		resp, err := HTCLIENT.Post(
			strings.TrimRight(addr, " /")+"/email/recv",
			"application/json",
			bytes.NewReader(serialize(Req{
				Recv: pbhash,
				Data: i,
			})),
		)
		if err != nil {
			break
		}
		if resp.ContentLength > int64(MAXESIZE) {
			break
		}
		err = json.NewDecoder(resp.Body).Decode(&servresp)
		resp.Body.Close()
		if err != nil {
			break
		}
		if servresp.Return != 0 {
			continue
		}
		pack := lc.Package(servresp.Result).Deserialize()
		if pack == nil {
			continue
		}
		pack = client.Decrypt(pack)
		if pack == nil {
			continue
		}
		err = DATABASE.SetEmail(user, pack)
		if err == nil {
			count++
		}
		if count == MAXCOUNT {
			break
		}
	}
}

func getTexts(email *Email) [2]string {
	head := strings.Split(email.Head, FSEPARAT)[0]
	body := strings.Split(email.Body, FSEPARAT)[0]
	return [2]string{
		head,
		body,
	}
}

func getFiles(email *Email) [][2]string {
	var list [][2]string
	name := strings.Split(email.Head, FSEPARAT)[1:]
	data := strings.Split(email.Body, FSEPARAT)[1:]
	for i := range name {
		list = append(list, [2]string{
			name[i],
			data[i],
		})
	}
	return list
}

func newEmail(sender, head, body string) *lc.Message {
	return lc.NewMessage([]byte(IS_EMAIL), serialize(Email{
		SenderName: sender,
		Head:       head,
		Body:       body,
	})).WithDiff(POWSDIFF)
}

func getName(user *User) string {
	if user == nil {
		return ""
	}
	return user.Name
}

func makeResult(retcod int, result string) (int, string) {
	return retcod, result
}
