package gopeer

import (
	"testing"
)

func TestConnectHashname(t *testing.T) {
	initConnects()
	defer clearConnects()
	var (
		errorIsExist bool
		hashes = [3]string{
			"qxn2+NJOo10RHmncVtA7fZ27Is678XPA1sZPlJFZGK4=",
			"LdQYHNvHfXDNXN+j92bk9wKSYUJ7sx+pphsb++TJnrw=",
			"fK/1gDOZHPU57MzghX5oBQLk2ETPKagmKwVr5UHV1gA=",
		}
	)
	for i, client := range Clients {
		for _, hash := range hashes {
			if hash == client.Hashname() {
				continue
			}
			if !client.InConnections(hash) {
				errorIsExist = true
				t.Errorf("client[%d].InConnections(%s)", i, hash)
			}
			conn := client.Connections[hash]
			if hash != conn.Hashname() {
				errorIsExist = true
				t.Errorf("hash != conn.Hashname()")
			}
		}
	}
	if !errorIsExist {
		t.Logf("connect.Hashname() success")
	}
}

func TestConnectPublic(t *testing.T) {
	initConnects()
	defer clearConnects()
	var (
		errorIsExist bool
		pubkeys = [3]string{
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
	for i, client := range Clients {
		for _, pub := range pubkeys {
			hash := HashPublic(ParsePublic(pub))
			if hash == client.Hashname() {
				continue
			}
			if !client.InConnections(hash) {
				errorIsExist = true
				t.Errorf("client[%d].InConnections(%s)", i, hash)
			}
			conn := client.Connections[hash]
			if pub != StringPublic(conn.Public()) {
				errorIsExist = true
				t.Errorf("pub != StringPublic(conn.Public())")
			}
		}
	}
	if !errorIsExist {
		t.Logf("connect.Public() success")
	}
}

func TestConnectThrow(t *testing.T) {
	initConnects()
	defer clearConnects()
	var (
		errorIsExist bool 
	)
	throw := Clients[0].Connections[Clients[2].Hashname()].Throw()
	if StringPublic(throw) != StringPublic(Clients[1].Public()) {
		errorIsExist = true
		t.Errorf("StringPublic(throw) != StringPublic(Clients[1].Public()) [1]")
	}
	throw = Clients[0].Connections[Clients[1].Hashname()].Throw()
	if StringPublic(throw) != StringPublic(Clients[1].Public()) {
		errorIsExist = true
		t.Errorf("StringPublic(throw) != StringPublic(Clients[1].Public()) [2]")
	}
	throw = Clients[2].Connections[Clients[0].Hashname()].Throw()
	if StringPublic(throw) != StringPublic(Clients[1].Public()) {
		errorIsExist = true
		t.Errorf("StringPublic(throw) != StringPublic(Clients[2].Public()) [1]")
	}
	throw = Clients[2].Connections[Clients[1].Hashname()].Throw()
	if StringPublic(throw) != StringPublic(Clients[1].Public()) {
		errorIsExist = true
		t.Errorf("StringPublic(throw) != StringPublic(Clients[2].Public()) [2]")
	}
	if !errorIsExist {
		t.Logf("connect.Throw() success")
	}
}
