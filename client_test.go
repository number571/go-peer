package gopeer

import (
	"testing"
	"time"
)

func TestClientHashname(t *testing.T) {
	var (
		errorIsExist bool
		hashes       = [3]string{
			"qxn2+NJOo10RHmncVtA7fZ27Is678XPA1sZPlJFZGK4=",
			"LdQYHNvHfXDNXN+j92bk9wKSYUJ7sx+pphsb++TJnrw=",
			"fK/1gDOZHPU57MzghX5oBQLk2ETPKagmKwVr5UHV1gA=",
		}
	)
	for index := range Clients {
		if Clients[index].Hashname() != hashes[index] {
			t.Errorf("client[%d].Hashname() != '%s'", index, hashes[index])
			errorIsExist = true
		}
	}
	if !errorIsExist {
		t.Logf("client.Hashname() success")
	}
}

func TestClientPublic(t *testing.T) {
	var (
		errorIsExist bool
		pubkeys      = [3]string{
			`-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAO3/0l65YYzobHcqJtlr/EkOul4N0RAKxhNOALfN00pRaApYpYVtgnWx
k9F4O75kXNEvh6ul84knyZF4XjZTp3LC+uTuptTI6ZBIRBMJPbJOkGjgkjeSTqkE
zWvHWzAe1nro3cJFpTgSAASPa00/2wmp2saXm1NrAJo5EYsLq3dpAgMBAAE=
-----END RSA PUBLIC KEY-----
`,
			`-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBANlaJjJxAkCrw4IMOL6N1BWsbKGNrRGnWii8oWMF1Ze+E/5lvwKxbsj4
XZvX+BUrzYNWeUnLhcwUhGawO26zUYVxf1mVTUHP9gqg9q8P1x+nGreRzhUrYyYq
SGa+bHwYVl1H6qyZD0HI9XdZoOWnBHMBsTYBMGxZdHPt9DBljrTpAgMBAAE=
-----END RSA PUBLIC KEY-----
`,
			`-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAKW4RsTSz1iioJU6oWlwRHoUM084lv7XxA3e2Ee9HYAZZu3iRF5tKUG3
YBR4HS9bGfoOh3FgUM8qQNZgA9wwI6Ylk5DzehHiC3YMGEn3pIgnPv+9bOUHMZXh
Npjj3L6l11Y69WiFgSkyepRbvWpU1qbeh1gzY2AAWao2ulmxqWGDAgMBAAE=
-----END RSA PUBLIC KEY-----
`,
		}
	)
	for index := range Clients {
		if StringPublic(Clients[index].Public()) != pubkeys[index] {
			t.Errorf("StringPublic(client[%d].Public()) != '%s'", index, pubkeys[index])
			errorIsExist = true
		}
	}
	if !errorIsExist {
		t.Logf("StringPublic(client.Public()) success")
	}
}

func TestClientInConnections(t *testing.T) {
	initConnects()
	defer clearConnects()
	var errorIsExist bool
	for i := range Clients {
		for j := range Clients {
			if Clients[i].Hashname() == Clients[j].Hashname() {
				continue
			}
			if !Clients[i].InConnections(Clients[j].Hashname()) {
				t.Errorf("client[%d].InConnections(client[%d].Hashname())", i, j)
				errorIsExist = true
			}
		}
	}
	if !errorIsExist {
		t.Logf("client.InConnections() success")
	}
}

func TestClientDisconnect(t *testing.T) {
	initConnects()
	defer clearConnects()
	var (
		err          error
		errorIsExist bool
	)
	dest := Clients[0].Destination(Clients[2].Hashname())
	err = Clients[0].Disconnect(dest)
	if err != nil {
		errorIsExist = true
		t.Errorf("client[0].Disconnect(2)")
	}
	time.Sleep(500 * time.Millisecond)
	dest = Clients[0].Destination(Clients[1].Hashname())
	err = Clients[0].Disconnect(dest)
	if err != nil {
		errorIsExist = true
		t.Errorf("client[0].Disconnect(1)")
	}
	dest = Clients[2].Destination(Clients[1].Hashname())
	err = Clients[2].Disconnect(dest)
	if err != nil {
		errorIsExist = true
		t.Errorf("client[2].Disconnect(1)")
	}
	if !errorIsExist {
		t.Logf("client.Disconnect() success")
	}
	time.Sleep(500 * time.Millisecond)
}

func TestClientConnect(t *testing.T) {
	defer clearConnects()
	var (
		err          error
		errorIsExist bool
	)
	dest := &Destination{
		Address:     Clients[1].Address(),
		Certificate: Clients[1].Certificate(),
		Public:      Clients[1].Public(),
	}
	err = Clients[0].Connect(dest)
	if err != nil {
		errorIsExist = true
		t.Errorf("client[0].Connect(1)")
	}
	err = Clients[2].Connect(dest)
	if err != nil {
		errorIsExist = true
		t.Errorf("client[2].Connect(1)")
	}
	dest = &Destination{
		Receiver: Clients[2].Public(),
	}
	err = Clients[0].Connect(dest)
	if err != nil {
		errorIsExist = true
		t.Errorf("client[0].Connect(2)")
	}
	if !errorIsExist {
		t.Logf("client.Connect() success")
	}
}

func TestClientSendTo(t *testing.T) {
	initConnects()
	defer clearConnects()
	var errorIsExist bool
	for i := range Clients {
		for j := range Clients {
			if Clients[i].Hashname() == Clients[j].Hashname() {
				continue
			}
			sendTestPackage(t, &errorIsExist, Clients[i], Clients[j])
		}
	}
	sendRedirectPackage(t, &errorIsExist)
	if !errorIsExist {
		t.Logf("client.SendTo() success")
	}
}

func sendTestPackage(t *testing.T, eie *bool, sender *Client, receiver *Client) {
	hash := receiver.Hashname()
	dest := sender.Destination(receiver.Hashname())
	sender.SendTo(dest, &Package{
		Head: Head{
			Title:  TITLE_TEST,
			Option: settings.OPTION_GET,
		},
		Body: Body{
			Data: "hello, world!",
		},
	})
	select {
	case <-sender.Connections[hash].Action:
		// pass
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		*eie = true
		t.Errorf("client[%s].SendTo(%s)", sender.Hashname(), receiver.Hashname())
	}
}

/*
      A -> F
==================
   A--------->B
             /|
   /--------- |
   |          |
   v          v
   C<-------->D
   |          |
   v          v
   \--------->E
              |
              v
              F
==================

qxn2+NJOo10RHmncVtA7fZ27Is678XPA1sZPlJFZGK4= ----- Clients[0] = B
LdQYHNvHfXDNXN+j92bk9wKSYUJ7sx+pphsb++TJnrw= ----- Clients[1] = C
fK/1gDOZHPU57MzghX5oBQLk2ETPKagmKwVr5UHV1gA= ----- Clients[2] = D

T3kBS/dw2afVaEhg6/vOHT7HH1PKgEdFjDWEN81NkJk= ----- newClient1 = A ----- SET
CoCV95aLEumtNXTEUHbXCN1GuCmISLpm5JcjROBJ914= ----- newClient2 = E
w38QEKtGkBJQoh1YBXK9w8pYNSmWe8MH4UP7jFxFnfw= ----- newClient3 = F ----- GET
*/

func sendRedirectPackage(t *testing.T, eie *bool) {
	var (
		newClient1 = new(Client) // A
		newClient2 = new(Client) // E
		newClient3 = new(Client) // F
	)
	defer func() {
		newClient1.listener.Close()
		newClient2.listener.Close()
		newClient3.listener.Close()
	}()
	var (
		privkeys = [3]string{
			`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDB2S1a/DQmkidrhpzlXvXV0YQWg71EaslFcXRmPBvjYIJ+TwyX
eTYvpsZMdJOV3H39lresRFhZewNiCoCeV71WZnnK4vJhPkqd0BlaGfDUeGHJq5/c
BrL2IuTCX0ZTK38XycbodoPeK6xaj8LS78ykgoVizzVy8cCjUmgI1NHl6QIDAQAB
AoGANQEtZbOQjvtny/8C57UPg2fGtmFPX2XToyliqpHFDmdVRzRWPRTnvB/eGQnH
UTL5QF312yTVA6KfSi+U+7cCDiP8gJR3lDK1Tw8cLAkiRguPyd778X9zdNSp5GGQ
Z89Y4Riedg0CSyjCpr8SkuZf7gWLbVOvY5dUHEvWEYfxmAECQQDqOtFw2Mm7bSjK
CNjRikXtS6eKmSpY0dQvpuxsDpjod/2qjcsjGzAMdtyD6aLo3OdJv/cfmxPBJki2
U/OFzSOzAkEA092JqRGyVe3jltF9F70ydb74UWZ1ZvOV8s16W5tANU25ajNrVVjK
w5VFVgA+waOP4uwCGT/7b86lH+4tkM1x8wJBAOVa1OzoCplRqUgz44NDH3fdxd/M
XQ/93wUOMaEZpha0MBrKn2fv3lvGI+WzaEcW0A+CPoyfQHe4cii/Cc0x80sCQBZo
1GrBqG9WXkBRoD2mkXPK41EY2Uoucang+hQ+c6gvtHD7R1sbrTbKzy6tj/XDazDB
bedl0R6eaPDbrI4obOkCQAVRnRwamYQUF3PDCEZNETy6lEn7P1Gu/gmjddMDzY/C
5Hr+uSH6GG/nel5PuWL05kfF/4Zo6bqAfk25Ylsshfc=
-----END RSA PRIVATE KEY-----
`,
			`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDI7fpjRFuFwSyHM61kHJjcESBNJxzMDOXPv9blmGof//FCv3l7
LOn+y7x7wrttFDNsUq0Tz8XtoFKcHZfnfaDSDvFOupOq9gvuE+chMPME/aiWDoP3
ul3IOuPzOvj2U7c6423VipIRpjRhMxTrDjgccC96+ZLUoHpjwAt0LZRjNwIDAQAB
AoGAag/EldWlggsgGBYdNrUpszCPHmrA9qzwOiOHqhY0HsQZdCRiWbSxp7+ftKxs
Bv4cztctaUHJts9kC+hLIPTdiFOA2J+TQA73UQisbnx9xLK7O7G9J9rmXjFUydnn
ZtaRbJzzwKS0N+R18VRpvquhyAJ6bROmmM4eCTLvY/EI4pECQQDMU9Sl540eVUA9
+msjOP7GOSwwYQvnM4WTY3Kk5JCe8sRDnAj7bIsiZUppFhYSTMPrYSLKdnE39wSx
jr+Uxza7AkEA+74rcaI3yS8ykHxeXr3KH7I0Naiox4A2UCiM6OSDUlgDXza1zsvN
SOAj9WlTmtdS3PfJ2XMOI4IWuHa0chuDtQJAIyfZUqJAeZOZdhB8Fpdb3qc/nyNT
dPg8Z1uZAh4BdBe5BRj3wqquVcCvyNzv9z4WB42y+rreSA7MU/CHrgWIuQJBAPGB
hvwEu6/t73xdU8tgF8BAnYW8v+5kOba8sDHcx37/oHx/Z/tz2QTIwrZ0zRgG6h/C
N4q8rhuyeUmN156AduECQQDDZC9mIzsuJOpmGWlbNKvvXaUIoGUp/Kyx/PbLDuzG
f3dwbf/3US7Lel+pLHqWNDThrkQouoved+6k+NgKaqph
-----END RSA PRIVATE KEY-----
`,
			`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCuvI9z+7squnL8JmeM364Bie7G7j3ZLPmyIeF11ydgGPwpyiOy
zAcKjdNQLYxt9Y2IMlguCNBQweHshnfhR5AeWY/X6F6SFWwyNqt5i0qV6ePUEln1
o56s+OOuO4fnAXyE6MZkYE54GSzdFzvJNK3LesvRtlsNXTbT/wzeg+fnlwIDAQAB
AoGAFAllPbyA8t5cbxOMTjgcAOsrKT6KcwvYOhfufY/FYRycVcJrI4aHzCsFLy15
6+X7a8GyIN073fbNjEzDFF8bZF7XIiNCcVmV/bAfcJ+i9ST4A+nSzrszkB+v0FMa
0JAl1Yh7rSylIsbYwM2zC6IDKLN8WFms0Qd4lHuASF2qTaECQQDBUKMq6qm5/STZ
rev9u1EOzN6lGFZs1LP4mFfX49tSJErMaKupRAyl1l8uKoVqzad+M7Ml2xooLQhA
+v/twVjRAkEA52WzUB4lD3Tvp44IOBq0juFgzkiWHkk6q5w97uUE8PmJXqZcN47i
iwd49J7lhzePx/Y/RDAlNIGsQ3FtSZVT5wJAOjNf3KTn0qIfPRY6zZperhkKEyR/
qKZlRLwA/nOQbWuVxXLh88UUFb2zzD9rCZu/CKTiE8yiVGQybvXipZ8ncQJAeMf8
8LTLY2YGMc9ROve1h17cyM/ai7Rti2Xibe/cxGt76IutVtKeLTOZTxYheJLn2dgO
7eizgtSstgdepCntwQJAHt3LiGbGQqbtwrPYkFbZPuYrwvGZOus/wGEfgrESODyv
RLohWeibe3Scpj0prVo5V1h2wmltUDu+ZoJWSObIpw==
-----END RSA PRIVATE KEY-----
`,
		}
	)
	newClient1 = createNewClient(":7070", ParsePrivate(privkeys[0]))
	dest := &Destination{
		Address:     newClient1.Address(),
		Certificate: newClient1.Certificate(),
		Public:      newClient1.Public(),
	}
	Clients[0].Connect(dest) // B -> A

	newClient2 = createNewClient(":9090", ParsePrivate(privkeys[1]))
	dest = &Destination{
		Address:     newClient2.Address(),
		Certificate: newClient2.Certificate(),
		Public:      newClient2.Public(),
	}
	Clients[1].Connect(dest) // C -> E
	Clients[2].Connect(dest) // D -> E

	newClient3 = createNewClient(settings.IS_CLIENT, ParsePrivate(privkeys[2]))
	newClient3.Connect(dest) // F -> E

	dest = &Destination{
		Receiver: newClient3.Public(),
	}

	err := newClient1.Connect(dest) // A -> F
	if err != nil {
		*eie = true
		t.Errorf("=================== client[A].Connect(F)")
		return
	}
	newClient1.SendTo(dest, &Package{
		Head: Head{
			Title:  TITLE_TEST,
			Option: settings.OPTION_GET,
		},
		Body: Body{
			Data: "hello, world!",
		},
	})

	select {
	case <-newClient1.Connections[newClient3.Hashname()].Action:
		// pass
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		*eie = true
		t.Errorf("=================== client[A].SendTo(F)[2]")
	}
}
