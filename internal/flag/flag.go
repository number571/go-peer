package flag

import (
	"os"
	"strings"
)

func GetFlagValue(pKey, pDefault string) string {
	return getFlagValueBySlice(os.Args[1:], pKey, pDefault)
}

func getFlagValueBySlice(args []string, pKey, pDefault string) string {
	isNextValue := false
	for _, arg := range args {
		if isNextValue {
			return arg
		}
		trimArg := strings.TrimLeft(arg, "-")
		if !strings.HasPrefix(trimArg, pKey) {
			continue
		}
		splited := strings.Split(trimArg, "=")
		if len(splited) == 1 {
			isNextValue = true
			continue
		}
		return strings.Join(splited[1:], "=")
	}
	if isNextValue {
		panic("args has key but value is not found")
	}
	return pDefault
}
