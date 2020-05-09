package gopeer

import (
	// "fmt"
	"crypto/rsa"
)

const (
	TITLE_TEST = "[TITLE-TEST]"
)

var (
	Clients [3]*Client
)

func init() {
	Set(SettingsType{
		"KEY_SIZE": uint64(1 << 10),
	})

	Clients[0] = createNewClient(settings.IS_CLIENT, ParsePrivate(PRIVATE_KEY1))
	Clients[1] = createNewClient(":8080", ParsePrivate(PRIVATE_KEY2))
	Clients[2] = createNewClient(settings.IS_CLIENT, ParsePrivate(PRIVATE_KEY3))
}

func createNewClient(address string, privkey *rsa.PrivateKey) *Client {
	var client = new(Client)
	nodeKey, nodeCert := GenerateCertificate(settings.NETWORK, settings.KEY_SIZE)
	listener := NewListener(address)
	listener.Open(&Certificate{
		Cert: []byte(nodeCert),
		Key:  []byte(nodeKey),
	}).Run(handleServer)
	if privkey == nil {
		client = listener.NewClient(GeneratePrivate(settings.KEY_SIZE))
	} else {
		client = listener.NewClient(privkey)
	}
	return client
}

func initConnects() {
	dest := &Destination{
		Address:     Clients[1].Address(),
		Certificate: Clients[1].Certificate(),
		Public:      Clients[1].Public(),
	}

	Clients[0].Connect(dest)
	Clients[2].Connect(dest)

	dest = &Destination{
		Receiver: Clients[2].Public(),
	}

	Clients[0].Connect(dest)
}

func clearConnects() {
	for i := range Clients {
		Clients[i].Action(func() {
			Clients[i].Connections = make(map[string]*Connect)
		})
	}
}

func handleServer(client *Client, pack *Package) {
	client.HandleAction(TITLE_TEST, pack,
		func(client *Client, pack *Package) (set string) {
			// fmt.Printf("[%s]: '%s'\n", pack.From.Sender.Hashname, pack.Body.Data)
			return set
		},
		func(client *Client, pack *Package) {
			client.Connections[pack.From.Sender.Hashname].Action <- true
		},
	)
}

const (
	PRIVATE_KEY1 = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDt/9JeuWGM6Gx3KibZa/xJDrpeDdEQCsYTTgC3zdNKUWgKWKWF
bYJ1sZPReDu+ZFzRL4erpfOJJ8mReF42U6dywvrk7qbUyOmQSEQTCT2yTpBo4JI3
kk6pBM1rx1swHtZ66N3CRaU4EgAEj2tNP9sJqdrGl5tTawCaORGLC6t3aQIDAQAB
AoGBALGT9mAtC8c6SGvlbJU/iD7umUnCH2JL15zhz5FVJrjF4s4NdHsIyZWNSNC7
WEBn3AVM5HrDWOHWaQR8fjck1cNPrPyGnAK/RxBSV6a2NDqnyFWa5JuXBZ84KvxR
jw3BVSrBlqKKmBiR/GvVKEAqW75hypZoUDgQ6YZZTh8I/Kb9AkEA85q1DOqddB0k
feTH0UMD1C7ou0JptJjUUuRYPbmDepLyNhKBu+SJRXlTMgKBeDNSGfuKEZyeu8A8
NphCTHeFBwJBAPocGeVNiQK11pEyqFQynF0prhLciqxUesJZITWWbr4vIFr1qD4Z
ekCCM7Bd/8oK/AwRmUmyalEU0vN2NZMg9A8CQGx/OiYPlKMzm54quEhepaTqY6OL
l9LkwqRMqXSMXJ/KNPCaW4fY6L61o7VBYnKrwORroPnpHNWYb/kM5XJzRR8CQQDJ
2Efl0G8UKt/hCjriyH18ihibzDR14y+3DOtKLf9tqOa5watngnQw/2LroNC/o6HJ
s6I74ar/iIi+RtXxyRRtAkAHlqSzMpn/wwoP504vuq/XpOYH7k/yMprlZeiMfXTq
iJO6tApJok6IemCsXan0KuuyrUiVLDKYEVL0Fr5mn4Zg
-----END RSA PRIVATE KEY-----
`

	PRIVATE_KEY2 = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZWiYycQJAq8OCDDi+jdQVrGyhja0Rp1oovKFjBdWXvhP+Zb8C
sW7I+F2b1/gVK82DVnlJy4XMFIRmsDtus1GFcX9ZlU1Bz/YKoPavD9cfpxq3kc4V
K2MmKkhmvmx8GFZdR+qsmQ9ByPV3WaDlpwRzAbE2ATBsWXRz7fQwZY606QIDAQAB
AoGAMS7zIcrsxZGYph754DTb3yOrWUtj9HC4OCunIW86jCHZCGIhvQxFk3iQRimy
26eT07qHq6lAV5P0+f+7EyvEsFCMx83ae95t8yS9WrwYpPFlJ1mg+cmfpCrRCgtd
5TtK4AnxhBfYWlTQwfdhjr/EUu5W+pYDATSwN0WFZas0fgECQQD/YMmk38aPNmBV
iAQj8uZeDF3+AcvMw2FtMHgJKk9kQuGe4mKgWHmQcdMih4dqZVnv3I8xsb3BZmmF
7xU719vJAkEA2eGnnLNUXMkqtt1CUkSHjU8+ShUTFiqy/C8Duf1JpiktEXFHFR67
kIh+9fWZLpYBgy0w9f8yNQOC5toRYdBgIQJAQ6WPxGzCXA07V2zALAWboC4Gd9Jh
+cuHczTzlvnuLdDJkxzEo1TMXsbH9s2PwU83k6IJDFDYwvIt4ZyDM2bqgQJBAL0j
KThbaCF/s+e4HMmDmdQudRkkQERe3q8SNP7whE2MswXQOu93lUT7aJMlF0uchkWU
Jkt1s+TXXnv901cA52ECQQDNqxcWi1MrLLTM2e3RaxUDQ5hOPCoeeLn19Jnyy+oK
+09atPvWw9facDI2o5BkCtCzomL5WYy3zk4DPF+fmMPN
-----END RSA PRIVATE KEY-----
`

	PRIVATE_KEY3 = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCluEbE0s9YoqCVOqFpcER6FDNPOJb+18QN3thHvR2AGWbt4kRe
bSlBt2AUeB0vWxn6DodxYFDPKkDWYAPcMCOmJZOQ83oR4gt2DBhJ96SIJz7/vWzl
BzGV4TaY49y+pddWOvVohYEpMnqUW71qVNam3odYM2NgAFmqNrpZsalhgwIDAQAB
AoGAKnLIVdQ13ltRzMfG3q4uLCMOfYVeUArOokrplX6gltZq4hYqMxG9FqM1Diky
OJYaUk30bZshN9932jGf63+9MYEPEvBcEH0nVLVE5Q/db66cf0Kh9pkwiShcASII
BLlL5fUexBi3XDPFPBwQVztcr/QPOe/NhYYwGSDUeDxkLYkCQQDCd9a/KpuP63mB
glNs6c367Q7Borj4bioWJHNlt5Q1cU79ueKqPSBZq8Agz3nMuce3P09zYoFE8hOt
dbhTLWvHAkEA2ifKQYBi283C22qScfyJlhNL9mNK853TJLURnyDeeXpeJeeF7Ocm
EE+V8e/D2qGWM6lb7Cj8zyycoyNWvPtEZQJAeUzk/6MlG5WG2fif7wy7tewOS0wj
0pps2BjufiEPanJ+EhfTwdqVBjnygsTHtaKgZ7Yu6csk1QumqIkIa6GmWwJADAsa
TVdrHbtUQIy3nPdWGSTjkqyUnLJfz6z3VhOYdJhezjTj3do87bWXD44u/8jf4+Y7
nuP8YOuTkiYHSdONSQJANMN0IRrsHdWsU3XIKrvrzXNuB+5wrFHy0gxgx32kInUm
vj51iwDgF7GxCdNVuRbDTyPcCYWqsiU8rDIdTUHXQQ==
-----END RSA PRIVATE KEY-----
`
)
