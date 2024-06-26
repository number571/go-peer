package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	initDir    = "_init/"
	keysDir    = "_keys/"
	startInput = "> "
	queueSize  = 32
)

var (
	queueVoid = make(chan []byte, queueSize)
	queue     = make(chan []byte, queueSize)
)

var (
	authBytes string
	connects  []string
	privKey   *rsa.PrivateKey
	attach    *rsa.PublicKey
)

func initApp() {
	if len(os.Args) != 3 {
		panic("example run: ./main [nickname] [host:port]")
	}

	authBytes = getAuthKey()
	connects = getConnects()
	privKey = getPrivKey()
	attach = &privKey.PublicKey
}

func main() {
	initApp()

	go runService()
	go runQueueVoid()
	go runQueue()

	for {
		cmd := readCmd(startInput)
		if len(cmd) != 2 {
			fmt.Println("len cmd != 2")
			continue
		}
		switch cmd[0] {
		case "attach":
			if err := getPubKey(cmd[1], attach); err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("ok")
		case "send":
			msg := fmt.Sprintf("%s%s: %s", authBytes, os.Args[1], strings.TrimSpace(cmd[1]))
			encBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, attach, []byte(msg), nil)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			queue <- encBytes
		}
	}
}

func runService() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		encBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer w.WriteHeader(http.StatusOK)
		msg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encBytes, nil)
		if err != nil {
			return
		}
		if !bytes.HasPrefix(msg, []byte(authBytes)) {
			return
		}
		msg = bytes.TrimPrefix(msg, []byte(authBytes))
		fmt.Printf("\n%s\n%s", string(msg), startInput)
	})
	_ = http.ListenAndServe(os.Args[2], nil) //nolint:gosec
}

func runQueue() {
	for {
		time.Sleep(5 * time.Second)
		encBytes := <-getQueue()
		for _, conn := range connects {
			go func(conn string) {
				req, err := http.NewRequestWithContext(
					context.Background(),
					http.MethodPost,
					conn,
					bytes.NewBuffer(encBytes),
				)
				if err != nil {
					panic(err)
				}
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
			}(conn)
		}
	}
}

func getQueue() chan []byte {
	if len(queue) == 0 {
		return queueVoid
	}
	return queue
}

func runQueueVoid() {
	for {
		if len(queueVoid) == queueSize {
			time.Sleep(time.Second)
			continue
		}
		encBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privKey.PublicKey, []byte("123"), nil)
		if err != nil {
			panic(err)
		}
		queueVoid <- encBytes
	}
}

func getPubKey(filename string, pubKey *rsa.PublicKey) error {
	pubKeyBytes, err := os.ReadFile(keysDir + filename)
	if err != nil {
		return err
	}
	pubKeyBlock, _ := pem.Decode(pubKeyBytes)
	if pubKeyBlock == nil || pubKeyBlock.Type != "PUBLIC KEY" {
		panic("public key: pem block is invalid")
	}
	pub, err := x509.ParsePKCS1PublicKey(pubKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}
	*pubKey = *pub
	return nil
}

func getAuthKey() string {
	authKeyBytes, err := os.ReadFile(initDir + "auth.key")
	if err != nil || len(authKeyBytes) == 0 {
		panic(err)
	}
	return string(authKeyBytes)
}

func getPrivKey() *rsa.PrivateKey {
	privKeyBytes, err := os.ReadFile(initDir + "priv.key")
	if err != nil {
		panic(err)
	}
	privateKeyBlock, _ := pem.Decode(privKeyBytes)
	if privateKeyBlock == nil || privateKeyBlock.Type != "PRIVATE KEY" {
		panic("private key: pem block is invalid")
	}
	priv, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}
	return priv
}

func getConnects() []string {
	cFile, err := os.Open(initDir + "connects.txt")
	if err != nil {
		panic(err)
	}
	defer cFile.Close()

	connects := make([]string, 0, 100)
	scanner := bufio.NewScanner(cFile)
	for scanner.Scan() {
		conn := strings.TrimSpace(scanner.Text())
		if conn == "" {
			continue
		}
		connects = append(connects, conn)
	}
	return connects
}

func readCmd(s string) []string {
	fmt.Print(s)
	input, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		panic(err)
	}
	cmd := strings.Split(string(input), "$")
	for i := range cmd {
		cmd[i] = strings.TrimSpace(cmd[i])
	}
	return cmd
}
