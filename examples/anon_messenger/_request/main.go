package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"

	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":"%s"
	}`

	cJsonDataTemplate = `{
        "method":"POST",
        "host":"go-peer/hidden-lake-messenger",
        "path":"/push",
        "head":{
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
		encryptMessage(
			getFriendPubKey(receiver),
			message,
		),
	)
}

func sendMessage(pReceiver string, pEncMessage []byte) {
	httpClient := http.Client{Timeout: time.Minute / 2}
	replacer := strings.NewReplacer("\n", "", "\t", "", "\r", "", " ", "", "\"", "\\\"")

	requestData := replacer.Replace(fmt.Sprintf(cJsonDataTemplate, base64.StdEncoding.EncodeToString(pEncMessage)))
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

func encryptMessage(pPubKey asymmetric.IPubKey, pMessage string) []byte {
	sessionKey := random.NewStdPRNG().GetBytes(32)
	message := getRandomMessageType(pMessage)
	return bytes.Join(
		[][]byte{
			[]byte(encoding.HexEncode(pPubKey.EncryptBytes(sessionKey))),
			[]byte(hlm_settings.CSeparator),
			symmetric.NewAESCipher(sessionKey).EncryptBytes(message),
		},
		[]byte{},
	)
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

func getFriendPubKey(pReceiver string) asymmetric.IPubKey {
	httpClient := http.Client{Timeout: time.Minute / 2}

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:7572/api/config/friends",
		nil,
	)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	vFriends := []hls_settings.SFriend{}
	if err := json.NewDecoder(resp.Body).Decode(&vFriends); err != nil {
		panic(err)
	}

	vPubKey := ""
	for _, friend := range vFriends {
		if friend.FAliasName == pReceiver {
			vPubKey = friend.FPublicKey
			break
		}
	}

	if vPubKey == "" {
		panic("friend undefined")
	}

	pubKey := asymmetric.LoadRSAPubKey(vPubKey)
	if pubKey == nil {
		panic("invalid public key")
	}

	return pubKey
}
