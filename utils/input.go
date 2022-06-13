package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InputString(begin string) string {
	fmt.Print(begin)
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(msg)
}
