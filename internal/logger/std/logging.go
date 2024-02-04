package std

import "fmt"

var (
	_ ILogging = sLogging{}
)

type sLogging []bool

func LoadLogging(pLogging []string) (ILogging, error) {
	mapping := map[string]int{
		CLogInfo: 0,
		CLogWarn: 1,
		CLogErro: 2,
	}

	logging := make(sLogging, len(mapping))

	for _, v := range pLogging {
		logType, ok := mapping[v]
		if !ok {
			return nil, fmt.Errorf("undefined log type '%s'", v)
		}
		logging[logType] = true
	}

	return logging, nil
}

func (p sLogging) HasInfo() bool {
	return p[0]
}

func (p sLogging) HasWarn() bool {
	return p[1]
}

func (p sLogging) HasErro() bool {
	return p[2]
}
