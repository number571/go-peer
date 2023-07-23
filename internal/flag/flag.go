package flag

import "flag"

func GetFlagValue(key, def, desc string) string {
	var (
		inputValue string
	)

	if flag.Lookup(key) == nil {
		flag.StringVar(&inputValue, key, def, desc)
		flag.Parse()
		return inputValue
	}

	return flag.Lookup(key).Value.(flag.Getter).Get().(string)
}
