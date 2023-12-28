package flag

import (
	"os"
	"strings"
)

func GetBoolFlagValue(pKey string) bool {
	return getBoolFlagValueBySlice(os.Args[1:], pKey)
}

func GetFlagValue(pKey, pDefault string) string {
	return getFlagValueBySlice(os.Args[1:], pKey, pDefault)
}

func getBoolFlagValueBySlice(args []string, pKey string) bool {
	for _, arg := range args {
		trimArg := strings.TrimLeft(arg, "-")
		if trimArg == pKey {
			return true
		}
	}
	return false
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
