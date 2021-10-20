package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Request struct {
	Host   string
	Path   string
	Method string
	Head   map[string]string
	Body   []byte
}

const (
	HLS                = "hidden-lake-service"
	FileWithPubKey     = "pub.key"
	ServerAddressInHLS = "route-service"
	AddressHLS         = "localhost:9571"
)

func fileIsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return data
}

func writeFile(file string, data []byte) error {
	return ioutil.WriteFile(file, data, 0644)
}

func serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func deserialize(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
