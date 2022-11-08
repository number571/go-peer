package inputter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var (
	_ IInputter = &sInputter{}
)

type sInputter struct {
	fBegin string
}

func NewInputter(begin string) IInputter {
	return &sInputter{
		fBegin: begin,
	}
}

func (inp *sInputter) String() string {
	fmt.Print(inp.fBegin)

	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(msg)
}

func (inp *sInputter) Password() string {
	fmt.Print(inp.fBegin)

	bpasw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	return strings.TrimSpace(string(bpasw))
}
