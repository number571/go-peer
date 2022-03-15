package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func InputString(begin string) string {
	fmt.Print(begin)
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(msg)
}

func ReadFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return data
}

func WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(file, data, 0644)
}

func Serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func Deserialize(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
