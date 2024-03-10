package app

import (
	"os"
)

func (p *sApp) initStorage() error {
	return os.MkdirAll(p.fStgPath, 0o777)
}
