package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/pkg/settings"
)

const (
	cSecretKey = "abc"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":"%s"
	}`

	cJsonDataTemplate = `{
        "method":"POST",
        "host":"hidden-lake-messenger",
        "path":"/push",
        "head":{
			"%s": "%s",
			"%s": "%s",
            "Accept": "application/json"
        },
        "body":"%s"
	}`
)

func main() {
	receiver := "Alice"
	message := "hello, world!"

	sendMessage(
		receiver,
		getRandomMessageType(message),
	)
}

func sendMessage(pReceiver string, pMessage []byte) {
	httpClient := http.Client{Timeout: time.Minute / 2}
	replacer := strings.NewReplacer("\n", "", "\t", "", "\r", "", " ", "", "\"", "\\\"")

	authKey := keybuilder.NewKeyBuilder(1, []byte(hlm_settings.CAuthSalt)).Build(cSecretKey)
	cipherKey := keybuilder.NewKeyBuilder(1, []byte(hlm_settings.CCipherSalt)).Build(cSecretKey)

	requestID := random.NewStdPRNG().GetString(hlm_settings.CRequestIDSize)
	pseudonym := "Bob"
	encMessage := symmetric.NewAESCipher(cipherKey).EncryptBytes(
		bytes.Join(
			[][]byte{
				hashing.NewHMACSHA256Hasher(
					authKey,
					bytes.Join([][]byte{[]byte(pseudonym), []byte(requestID), pMessage}, []byte{}),
				).ToBytes(),
				pMessage,
			},
			[]byte{},
		),
	)

	// fmt.Println(requestID)
	// for _, c := range encMessage {
	// 	fmt.Printf("\\x%02x", c)
	// }
	// fmt.Println()

	requestData := replacer.Replace(
		fmt.Sprintf(
			cJsonDataTemplate,
			hlm_settings.CHeaderPseudonym,
			pseudonym,
			hlm_settings.CHeaderRequestId,
			requestID,
			base64.StdEncoding.EncodeToString(encMessage),
		),
	)
	req, err := http.NewRequest(
		http.MethodPut,
		"http://localhost:7572/api/network/request",
		bytes.NewBufferString(fmt.Sprintf(cRequestTemplate, pReceiver, requestData)),
	)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}

func getRandomMessageType(pMessage string) []byte {
	if random.NewStdPRNG().GetBool() { // isText
		return bytes.Join(
			[][]byte{
				{hlm_settings.CIsText},
				[]byte(pMessage),
			},
			[]byte{},
		)
	}
	// isFile
	return bytes.Join(
		[][]byte{
			{hlm_settings.CIsFile},
			[]byte("example.txt"),
			{hlm_settings.CIsFile},
			[]byte(pMessage),
		},
		[]byte{},
	)
}
