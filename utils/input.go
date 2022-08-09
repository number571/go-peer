package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var (
	_ IInput = &sInput{}
)

type sInput struct {
	fBegin string
}

func NewInput(begin string) IInput {
	return &sInput{
		fBegin: begin,
	}
}

func (inp *sInput) String() string {
	fmt.Print(inp.fBegin)

	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(msg)
}

func (inp *sInput) Password() string {
	fmt.Print(inp.fBegin)

	bpasw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	return string(bpasw)
}
