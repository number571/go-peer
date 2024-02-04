package app

import (
	"os"
)

func (p *sApp) initStorage() error {
	return os.MkdirAll(p.fPathTo, 0o777)
}
