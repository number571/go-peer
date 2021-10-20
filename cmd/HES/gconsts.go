package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/url"
	"time"

	gp "github.com/number571/gopeer"
	"golang.org/x/net/proxy"
)

const (
	TMESSAGE = "\005\007\001\000\001\007\005"
)

var (
	MAXESIZE = uint(0)
	POWSDIFF = uint(0)
	OPENADDR = ""
	HTCLIENT = new(http.Client)
)

func hesDefaultInit(address string) {
	socks5Ptr := flag.String("socks5", "", "enable socks5 and create proxy connection")
	addrPtr := flag.String("open", address, "open address for hidden email server")
	flag.Parse()
	OPENADDR = *addrPtr
	if *socks5Ptr != "" {
		socks5, err := url.Parse("socks5://" + *socks5Ptr)
		if err != nil {
			panic("error: socks5 conn")
		}
		dialer, err := proxy.FromURL(socks5, proxy.Direct)
		if err != nil {
			panic("error: dialer")
		}
		HTCLIENT = &http.Client{
			Transport: &http.Transport{Dial: dialer.Dial},
			Timeout:   time.Second * 15,
		}
	}
	gp.Set(gp.SettingsType{
		"POWS_DIFF": uint(25),
		"PACK_SIZE": uint(8 << 20),
		"AKEY_SIZE": uint(2 << 10),
		"SKEY_SIZE": uint(1 << 5),
		"RAND_SIZE": uint(1 << 4),
	})
	MAXESIZE = gp.Get("PACK_SIZE").(uint)
	POWSDIFF = gp.Get("POWS_DIFF").(uint)
}

func serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}
