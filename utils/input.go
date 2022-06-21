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

func InputString(begin string) string {
	fmt.Print(begin)

	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(msg)
}

// Settings must contain (CSizePasw, CMaskPasw).
func InputPassword(sett settings.ISettings, begin string) string {
	fmt.Print(begin)

	bpsw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	spsw := strings.TrimSpace(string(bpsw))
	if uint64(len([]rune(spsw))) < sett.Get(settings.CSizePasw) {
		panic("length of password < min size")
	}

	maskPasw := sett.Get(settings.CMaskPasw)
	if !passwordHasMask(maskPasw, spsw) {
		panic(fmt.Sprintf("password mode should be = %03b", maskPasw))
	}

	return spsw
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
