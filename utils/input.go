package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unicode"

	"github.com/number571/go-peer/settings"
	"golang.org/x/term"
)

var (
	_ IInput = &sInput{}
)

type sInput struct {
	fBegin string
	fSett  settings.ISettings
}

func NewInput(sett settings.ISettings, begin string) IInput {
	return &sInput{
		fBegin: begin,
		fSett:  sett,
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

// Settings must contain (CSizePasw, CMaskPasw)
func (inp *sInput) Password() string {
	fmt.Print(inp.fBegin)

	bpasw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	spasw := strings.TrimSpace(string(bpasw))
	if uint64(len([]rune(spasw))) < inp.fSett.Get(settings.CSizePasw) {
		panic("length of password < min size")
	}

	maskPasw := inp.fSett.Get(settings.CMaskPasw)
	if !passwordHasMask(maskPasw, spasw) {
		panic(fmt.Sprintf("password mode should be = %03b", maskPasw))
	}

	return spasw
}

func passwordHasMask(maskPasw uint64, pasw string) bool {
	var (
		isAlphabet = maskPasw&settings.CPaswAplh == 0
		isNumeric  = maskPasw&settings.CPaswNumr == 0
		isSpecial  = maskPasw&settings.CPaswSpec == 0
	)

	for _, ch := range pasw {
		if unicode.IsLetter(ch) {
			isAlphabet = true
			continue
		}
		if _, err := strconv.Atoi(string(ch)); err == nil {
			isNumeric = true
			continue
		}
		isSpecial = true
	}

	return isAlphabet && isNumeric && isSpecial
}
